package core

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

var Http *HttpEngine

// HttpMetrics 用于记录HTTP性能指标
type HttpMetrics struct {
	requestCount      int64 // 总请求数
	errorCount        int64 // 错误请求数
	totalResponseTime int64 // 总响应时间（毫秒）
	activeRequests    int32 // 当前活跃请求数
}

// HttpEngineConfig HTTP引擎配置
type HttpEngineConfig struct {
	MaxConns           int           // 最大并发连接数
	PoolSize           int           // 工作池大小
	MaxRetries         int           // 最大重试次数
	ProxyURLs          []string      // 多个代理服务器URL
	ProxyTimeout       time.Duration // 代理超时时间
	ProxyRetries       int           // 代理重试次数
	Timeout            time.Duration // 请求超时时间
	RetryInterval      time.Duration // 重试间隔时间
	CompressionEnabled bool          // 是否启用压缩
}

// HttpRequest 请求结构体
type HttpRequest struct {
	ID         int64              // 请求唯一标识
	Method     string             // HTTP方法
	URL        string             // 请求URL
	Headers    map[string]string  // 请求头
	done       bool               // 请求是否完成
	Body       []byte             // 请求体
	Response   *HttpResponse      // 响应对象
	Err        error              // 错误信息
	Wg         sync.WaitGroup     // 等待组
	Callback   func(*HttpRequest) // 回调函数
	retries    int                // 已重试次数
	compressed bool               // 是否已压缩
}

// HttpResponse 响应结构体
type HttpResponse struct {
	raw     *http.Response // 原始响应
	code    int            // 状态码
	Headers http.Header    // 响应头
	Body    []byte         // 响应体
}

// HttpEngine HTTP引擎核心结构体
type HttpEngine struct {
	client          *http.Client                      // HTTP客户端
	sem             *semaphore.Weighted               // 信号量，用于限制并发
	maxRetries      int                               // 最大重试次数
	poolSize        int                               // 工作池大小
	tasks           chan *HttpRequest                 // 任务通道
	wg              sync.WaitGroup                    // 等待组
	logger          *logrus.Logger                    // 日志记录器
	metrics         *HttpMetrics                      // 性能指标
	config          *HttpEngineConfig                 // 配置信息
	patterns        []map[string]string               // 流量混淆模式
	errorHandlers   map[int]func(*HttpResponse) error // 错误处理器映射
	retryableCodes  map[int]bool                      // 可重试的状态码
	chunkConfig     *ChunkConfig                      // 分块配置
	transferMetrics *TransferMetrics                  // 传输指标
	cacheManager    *CacheManager                     // 缓存管理器
	metricsChan     chan metricsData                  // 指标数据通道
	cacheChan       chan *HttpRequest                 // 缓存请求通道
}

// 通用头部模板
var commonHeaders = []map[string]string{
	{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
		"Connection":      "keep-alive",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	},
	{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Connection":      "keep-alive",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	},
}

// 添加错误类型定义
type HttpError struct {
	Code    int
	Message string
	Err     error
}

func (e *HttpError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("HTTP Error %d: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("HTTP Error %d: %s", e.Code, e.Message)
}

// 定义错误码常量
const (
	// 4xx Client Errors
	ErrBadRequest       = 400
	ErrUnauthorized     = 401
	ErrForbidden        = 403
	ErrNotFound         = 404
	ErrMethodNotAllowed = 405
	ErrRequestTimeout   = 408
	ErrConflict         = 409
	ErrTooManyRequests  = 429

	// 5xx Server Errors
	ErrServerError        = 500
	ErrBadGateway         = 502
	ErrServiceUnavailable = 503
	ErrGatewayTimeout     = 504
)

// 用于传递指标数据的结构
type metricsData struct {
	duration   time.Duration
	statusCode int
}

func NewHttpEngine(config *HttpEngineConfig) *HttpEngine {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxConnsPerHost:     config.MaxConns,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		DisableCompression:  !config.CompressionEnabled,
		MaxIdleConnsPerHost: config.MaxConns / 2,
	}

	// 设置代理选择器
	basicCfg := GetInstance()
	if basicCfg.Proxy.Enabled && len(basicCfg.Proxy.ProxyPool) > 0 {
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			proxyURL := basicCfg.Proxy.GetProxyURL(req.URL.Scheme)
			if proxyURL == "" {
				return nil, nil
			}
			return url.Parse(proxyURL)
		}
	}

	// 初始化缓存
	cacheManager := GetCacheManager()
	if cacheManager == nil {
		logger.Error("Failed to initialize cache manager")
	}

	engine := &HttpEngine{
		client: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
		sem:           semaphore.NewWeighted(int64(config.MaxConns)),
		maxRetries:    config.MaxRetries,
		poolSize:      config.PoolSize,
		tasks:         make(chan *HttpRequest, config.PoolSize),
		logger:        logger,
		metrics:       &HttpMetrics{},
		config:        config,
		patterns:      commonHeaders,
		errorHandlers: make(map[int]func(*HttpResponse) error),
		retryableCodes: map[int]bool{
			408: true, // Request Timeout
			429: true, // Too Many Requests
			500: true, // Internal Server Error
			502: true, // Bad Gateway
			503: true, // Service Unavailable
			504: true, // Gateway Timeout
		},
		cacheManager: cacheManager,
		metricsChan:  make(chan metricsData, 1000), // 缓冲通道
		cacheChan:    make(chan *HttpRequest, 1000),
	}

	// 注册默认错误处理器
	engine.registerDefaultErrorHandlers()

	engine.initWorkerPool()

	// 启动异步处理协程
	go engine.processMetrics()
	go engine.processCache()

	return engine
}

func GetHttpEngine() *HttpEngine {
	if Http == nil {
		basicCfg := GetInstance()
		config := &HttpEngineConfig{
			MaxConns:           200,
			PoolSize:           10,
			MaxRetries:         3,
			ProxyURLs:          basicCfg.Proxy.ProxyPool,
			ProxyTimeout:       time.Duration(basicCfg.Proxy.Timeout) * time.Second,
			ProxyRetries:       basicCfg.Proxy.Retries,
			Timeout:            time.Duration(basicCfg.Timeout.Read) * time.Second,
			RetryInterval:      time.Second,
			CompressionEnabled: true,
		}
		Http = NewHttpEngine(config)
	}
	return Http
}

// 添加流量混淆
func (engine *HttpEngine) obfuscateRequest(req *HttpRequest) {
	// 随机选择一个请求头模板
	pattern := engine.patterns[time.Now().UnixNano()%int64(len(engine.patterns))]

	// 应用基础请求头
	for k, v := range pattern {
		if _, exists := req.Headers[k]; !exists {
			req.Headers[k] = v
		}
	}

	// 添加随机参数
	if req.Method == "GET" {
		separator := "?"
		if bytes.Contains([]byte(req.URL), []byte("?")) {
			separator = "&"
		}
		noise := make([]byte, 4)
		rand.Read(noise)
		req.URL += fmt.Sprintf("%s_=%s", separator, base64.URLEncoding.EncodeToString(noise))
	}

	// 添加随机延迟
	//time.Sleep(time.Duration(50+rand.Int63n(200)) * time.Millisecond)
}

// compressBody 压缩请求体
func compressBody(body []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(body)
	if err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}

// decompressBody 解压响应体
func decompressBody(body []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// ExecuteRequest 执行HTTP请求，支持重试机制
func (engine *HttpEngine) ExecuteRequest(req *HttpRequest) error {
	var lastErr error
	backoff := engine.config.RetryInterval // 初始重试间隔

	// 重试循环
	for attempt := 0; attempt <= engine.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * 1.5) // 指数退避策略
		}

		err := engine.executeRequestOnce(req)
		if err == nil {
			return nil
		}

		// 检查错误是否可重试
		if httpErr, ok := err.(*HttpError); ok {
			if !engine.retryableCodes[httpErr.Code] {
				return err // 不可重试的错误直接返回
			}
		}

		lastErr = err
		req.retries++

		engine.logger.Debugf("Request %d retry %d/%d: %v",
			req.ID, attempt+1, engine.maxRetries, err)
	}

	return fmt.Errorf("max retries exceeded: %v", lastErr)
}

func (engine *HttpEngine) executeRequestOnce(req *HttpRequest) error {
	ctx := context.Background()
	if err := engine.sem.Acquire(ctx, 1); err != nil {
		return &HttpError{Code: 0, Message: "Failed to acquire semaphore", Err: err}
	}
	defer engine.sem.Release(1)

	atomic.AddInt32(&engine.metrics.activeRequests, 1)
	defer atomic.AddInt32(&engine.metrics.activeRequests, -1)

	start := time.Now()

	// 流量混淆
	engine.obfuscateRequest(req)

	// 压缩大请求
	//if engine.config.CompressionEnabled && len(req.Body) > 1024 {
	//	compressed, err := compressBody(req.Body)
	//	if err == nil && len(compressed) < len(req.Body) {
	//		req.Body = compressed
	//		req.Headers["Content-Encoding"] = "gzip"
	//		req.compressed = true
	//	}
	//}

	// 构建请求
	httpReq, err := http.NewRequest(req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		return &HttpError{Code: 0, Message: "Failed to create request", Err: err}
	}

	// 设置请求头
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// 发送请求
	resp, err := engine.client.Do(httpReq)
	if err != nil {
		// 异步记录错误指标
		go func() {
			engine.metricsChan <- metricsData{
				duration:   time.Since(start),
				statusCode: 0,
			}
		}()
		return &HttpError{Code: 0, Message: "Request failed", Err: err}
	}

	if resp == nil {
		return &HttpError{Code: 0, Message: "Empty response", Err: nil}
	}

	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &HttpError{Code: resp.StatusCode, Message: "Failed to read response body", Err: err}
	}

	// 处理压缩的响应
	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, err = decompressBody(body)
		if err != nil {
			return &HttpError{Code: resp.StatusCode, Message: "Failed to decompress response", Err: err}
		}
	}

	req.Response = &HttpResponse{
		raw:     resp,
		code:    resp.StatusCode,
		Headers: resp.Header,
		Body:    body,
	}

	// 处理错误状态码
	if resp.StatusCode >= 400 {
		atomic.AddInt64(&engine.metrics.errorCount, 1)

		if handler, exists := engine.errorHandlers[resp.StatusCode]; exists {
			return handler(req.Response)
		}

		// 根据状态码范围返回通用错误
		switch {
		case resp.StatusCode >= 400 && resp.StatusCode < 500:
			return &HttpError{
				Code:    resp.StatusCode,
				Message: fmt.Sprintf("Client Error: %s", resp.Status),
				Err:     fmt.Errorf("unexpected 4xx error: %d", resp.StatusCode),
			}
		case resp.StatusCode >= 500:
			return &HttpError{
				Code:    resp.StatusCode,
				Message: fmt.Sprintf("Server Error: %s", resp.Status),
				Err:     fmt.Errorf("unexpected 5xx error: %d", resp.StatusCode),
			}
		}
	}

	req.done = true

	// 更新指标
	duration := time.Since(start)
	atomic.AddInt64(&engine.metrics.requestCount, 1)
	atomic.AddInt64(&engine.metrics.totalResponseTime, duration.Milliseconds())

	if req.done {
		// 异步发送性能指标
		go func() {
			engine.metricsChan <- metricsData{
				duration:   duration,
				statusCode: resp.StatusCode,
			}
		}()

		// 异步处理缓存
		go func() {
			engine.cacheChan <- req
		}()
	}

	return nil
}

// 异步处理性能指标
func (engine *HttpEngine) processMetrics() {
	for data := range engine.metricsChan {
		atomic.AddInt64(&engine.metrics.requestCount, 1)
		atomic.AddInt64(&engine.metrics.totalResponseTime, data.duration.Milliseconds())
		if data.statusCode >= 400 {
			atomic.AddInt64(&engine.metrics.errorCount, 1)
		}
	}
}

// 异步处理缓存
func (engine *HttpEngine) processCache() {
	for req := range engine.cacheChan {
		if err := engine.cacheManager.SaveToCache(req); err != nil {
			engine.logger.Warnf("Failed to cache request %d: %v", req.ID, err)
		}
	}
}

// StopAndWait 优雅关闭HTTP引擎
func (engine *HttpEngine) StopAndWait() {
	// 关闭所有通道
	close(engine.metricsChan)
	close(engine.cacheChan)

	// 等待所有工作协程完成
	engine.wg.Wait()

	// 清理缓存
	if err := engine.cacheManager.CloseAndCleanup(); err != nil {
		engine.logger.Errorf("Failed to cleanup cache: %v", err)
	}
}

func (engine *HttpEngine) registerDefaultErrorHandlers() {
	// 4xx 错误处理
	engine.errorHandlers[ErrBadRequest] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrBadRequest,
			Message: "Bad Request: The request could not be understood or was missing required parameters",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrUnauthorized] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrUnauthorized,
			Message: "Unauthorized: Authentication is required and has failed or not been provided",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrForbidden] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrForbidden,
			Message: "Forbidden: You don't have permission to access this resource",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrNotFound] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrNotFound,
			Message: "Not Found: The requested resource could not be found",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrMethodNotAllowed] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrMethodNotAllowed,
			Message: "Method Not Allowed: The request method is not supported for this resource",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrRequestTimeout] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrRequestTimeout,
			Message: "Request Timeout: The server timed out waiting for the request",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrConflict] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrConflict,
			Message: "Conflict: The request conflicts with the current state of the server",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrTooManyRequests] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrTooManyRequests,
			Message: "Too Many Requests: You have exceeded the rate limit",
			Err:     nil,
		}
	}

	// 5xx 错误处理
	engine.errorHandlers[ErrServerError] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrServerError,
			Message: "Internal Server Error: An unexpected condition was encountered",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrBadGateway] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrBadGateway,
			Message: "Bad Gateway: The server received an invalid response from the upstream server",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrServiceUnavailable] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrServiceUnavailable,
			Message: "Service Unavailable: The server is temporarily unable to handle the request",
			Err:     nil,
		}
	}

	engine.errorHandlers[ErrGatewayTimeout] = func(resp *HttpResponse) error {
		return &HttpError{
			Code:    ErrGatewayTimeout,
			Message: "Gateway Timeout: The upstream server failed to respond in time",
			Err:     nil,
		}
	}
}

// 初始化可重试状态码
func (engine *HttpEngine) initRetryableCodes() {
	engine.retryableCodes = map[int]bool{
		ErrRequestTimeout:     true,
		ErrTooManyRequests:    true,
		ErrServerError:        true,
		ErrBadGateway:         true,
		ErrServiceUnavailable: true,
		ErrGatewayTimeout:     true,
	}
}

// IsClientError 判断是否为客户端错误（4xx）
func (engine *HttpEngine) IsClientError(err error) bool {
	if httpErr, ok := err.(*HttpError); ok {
		return httpErr.Code >= 400 && httpErr.Code < 500
	}
	return false
}

// IsServerError 判断是否为服务器错误（5xx）
func (engine *HttpEngine) IsServerError(err error) bool {
	if httpErr, ok := err.(*HttpError); ok {
		return httpErr.Code >= 500 && httpErr.Code < 600
	}
	return false
}

// IsRetryableError 判断错误是否可重试
func (engine *HttpEngine) IsRetryableError(err error) bool {
	if httpErr, ok := err.(*HttpError); ok {
		return engine.retryableCodes[httpErr.Code]
	}
	return false
}

// 分块输相关配置
type ChunkConfig struct {
	ChunkSize   int64         // 分块大小
	Concurrency int           // 并发数
	MaxRetries  int           // 最大重试次数
	RetryDelay  time.Duration // 重试延迟
}

// 分块传输请求
type ChunkRequest struct {
	*HttpRequest
	ChunkIndex  int   // 当前块索引
	TotalChunks int   // 总块数
	Offset      int64 // 文件偏移量
	Size        int64 // 块大小
}

// 传输性能指标
type TransferMetrics struct {
	BytesSent       int64        // 已发送字节数
	BytesReceived   int64        // 已接收字节数
	ChunksProcessed int64        // 已处理块数
	FailedChunks    int64        // 失败块数
	StartTime       time.Time    // 开始时间
	mu              sync.RWMutex // 指标锁
}

// 扩展 HttpEngine
func (engine *HttpEngine) InitChunkTransfer(config ChunkConfig) {
	engine.chunkConfig = &config
	engine.transferMetrics = &TransferMetrics{
		StartTime: time.Now(),
	}
}

// 处理分块请求
func (engine *HttpEngine) SendChunkRequest(chunk *ChunkRequest) error {
	// 设置分块传输的头部信息
	chunk.Headers["X-Chunk-Index"] = fmt.Sprintf("%d", chunk.ChunkIndex)
	chunk.Headers["X-Total-Chunks"] = fmt.Sprintf("%d", chunk.TotalChunks)
	chunk.Headers["X-Chunk-Offset"] = fmt.Sprintf("%d", chunk.Offset)
	chunk.Headers["X-Chunk-Size"] = fmt.Sprintf("%d", chunk.Size)
	chunk.Headers["Content-Length"] = fmt.Sprintf("%d", len(chunk.Body))

	// 发送请求并处理重试
	var err error
	for retry := 0; retry <= engine.chunkConfig.MaxRetries; retry++ {
		err = engine.ExecuteRequest(chunk.HttpRequest)
		if err == nil {
			engine.updateTransferMetrics(int64(len(chunk.Body)), int64(len(chunk.Response.Body)))
			return nil
		}

		if retry < engine.chunkConfig.MaxRetries {
			time.Sleep(engine.chunkConfig.RetryDelay * time.Duration(retry+1))
		}
	}

	engine.incrementFailedChunks()
	return fmt.Errorf("chunk transfer failed after %d retries: %v", engine.chunkConfig.MaxRetries, err)
}

// 并发处理多个分块请求
func (engine *HttpEngine) ProcessChunks(chunks []*ChunkRequest) error {
	sem := semaphore.NewWeighted(int64(engine.chunkConfig.Concurrency))
	var wg sync.WaitGroup
	errChan := make(chan error, len(chunks))

	for _, chunk := range chunks {
		wg.Add(1)
		go func(c *ChunkRequest) {
			defer wg.Done()

			// 获取信号量
			if err := sem.Acquire(context.Background(), 1); err != nil {
				errChan <- fmt.Errorf("failed to acquire semaphore: %v", err)
				return
			}
			defer sem.Release(1)

			// 发送分块请求
			if err := engine.SendChunkRequest(c); err != nil {
				errChan <- err
			}
		}(chunk)
	}

	// 等待所有分块处理完成
	wg.Wait()
	close(errChan)

	// 收集错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("chunk processing errors: %v", errors)
	}

	return nil
}

// 更新传输指标
func (engine *HttpEngine) updateTransferMetrics(sent, received int64) {
	engine.transferMetrics.mu.Lock()
	defer engine.transferMetrics.mu.Unlock()

	atomic.AddInt64(&engine.transferMetrics.BytesSent, sent)
	atomic.AddInt64(&engine.transferMetrics.BytesReceived, received)
	atomic.AddInt64(&engine.transferMetrics.ChunksProcessed, 1)
}

func (engine *HttpEngine) incrementFailedChunks() {
	atomic.AddInt64(&engine.transferMetrics.FailedChunks, 1)
}

// 获取传输状态
func (engine *HttpEngine) GetTransferStatus() map[string]interface{} {
	engine.transferMetrics.mu.RLock()
	defer engine.transferMetrics.mu.RUnlock()

	duration := time.Since(engine.transferMetrics.StartTime).Seconds()
	bytesSent := atomic.LoadInt64(&engine.transferMetrics.BytesSent)
	bytesReceived := atomic.LoadInt64(&engine.transferMetrics.BytesReceived)

	return map[string]interface{}{
		"bytes_sent":       bytesSent,
		"bytes_received":   bytesReceived,
		"chunks_processed": atomic.LoadInt64(&engine.transferMetrics.ChunksProcessed),
		"failed_chunks":    atomic.LoadInt64(&engine.transferMetrics.FailedChunks),
		"upload_speed":     float64(bytesSent) / duration,
		"download_speed":   float64(bytesReceived) / duration,
		"duration_seconds": duration,
	}
}

// 优化连接池配置
func (engine *HttpEngine) optimizeConnectionPool() {
	engine.client.Transport.(*http.Transport).MaxIdleConns = engine.chunkConfig.Concurrency * 2
	engine.client.Transport.(*http.Transport).MaxIdleConnsPerHost = engine.chunkConfig.Concurrency
	engine.client.Transport.(*http.Transport).IdleConnTimeout = 30 * time.Second
}

// InitDefaultChunkTransfer 使用默认配置初始化分块传输
func (engine *HttpEngine) InitDefaultChunkTransfer() {
	defaultConfig := ChunkConfig{
		ChunkSize:   1024 * 1024, // 1MB
		Concurrency: 3,
		MaxRetries:  3,
		RetryDelay:  time.Second,
	}
	engine.InitChunkTransfer(defaultConfig)
}

// AutoInitChunkTransfer 动检测是否需要使用分块传输
func (engine *HttpEngine) AutoInitChunkTransfer(contentLength int64) bool {
	// 定义触发分块传输的阈值（5MB）
	const chunkThreshold int64 = 5 * 1024 * 1024

	if contentLength > chunkThreshold {
		// 内容较大，初始化分块传输
		defaultConfig := ChunkConfig{
			ChunkSize:   1024 * 1024, // 1MB
			Concurrency: 3,
			MaxRetries:  3,
			RetryDelay:  time.Second,
		}
		engine.InitChunkTransfer(defaultConfig)
		return true
	}

	// 内容较小，不需要分块
	engine.chunkConfig = nil
	engine.transferMetrics = nil
	return false
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{}
}

// 添加缓存查询方法
func (engine *HttpEngine) GetCachedRequest(requestID int64) (*HttpRequest, error) {
	return engine.cacheManager.GetFromCache(requestID)
}

func (engine *HttpEngine) initWorkerPool() {
	engine.wg.Add(engine.poolSize)
	for i := 0; i < engine.poolSize; i++ {
		go func() {
			defer engine.wg.Done()
			for req := range engine.tasks {
				if err := engine.ExecuteRequest(req); err != nil {
					engine.logger.Errorf("Request %d failed: %v", req.ID, err)
				}
				if req.Callback != nil {
					req.Callback(req)
				}
			}
		}()
	}
}
