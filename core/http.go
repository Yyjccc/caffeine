package core

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var Http *HttpEngine

//封装http引擎，生产者消费者模型

// HttpRequest 结构体，表示一个HTTP请求
type HttpRequest struct {
	ID       int64              // 请求的唯一标识符
	Method   string             // 请求方法
	URL      string             // 请求的URL
	Headers  map[string]string  // 自定义请求头
	done     bool               // 请求是否已经完成
	Body     []byte             // 请求体（适用于POST等方法）
	Response *HttpResponse      //响应
	Err      error              //错误
	Wg       sync.WaitGroup     // 用于等待请求完成
	Callback func(*HttpRequest) // 请求完成后的回调函数
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{
		ID:      GenerateID(),
		Headers: make(map[string]string),
		done:    false,
		Wg:      sync.WaitGroup{},
	}
}

// HttpResponse 结构体，封装HTTP响应，方便解析
type HttpResponse struct {
	raw     *http.Response //原始响应
	code    int            // HTTP状态码
	Headers http.Header    // 响应头
	Body    []byte         // 响应体原始数据
}

// HttpEngine 定义HTTP引擎结构体
type HttpEngine struct {
	client     *http.Client        // HTTP客户端
	sem        *semaphore.Weighted // 信号量，用于控制最大并发
	maxRetries int                 // 最大重试次数
	poolSize   int                 // 协程池大小
	tasks      chan *HttpRequest   // 请求任务通道
	wg         sync.WaitGroup      // 用于等待所有任务完成
	logger     *logrus.Logger
}

// NewHttpEngine 创建新的HTTP引擎
func NewHttpEngine(maxConns, poolSize, maxRetries int) *HttpEngine {
	transport := &http.Transport{
		MaxIdleConns:      100,
		MaxConnsPerHost:   maxConns,
		IdleConnTimeout:   90 * time.Second,
		DisableKeepAlives: false,
	}

	engine := &HttpEngine{
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: transport,
		},
		sem:        semaphore.NewWeighted(int64(maxConns)), // 信号量控制最大并发数
		maxRetries: maxRetries,
		poolSize:   poolSize,
		tasks:      make(chan *HttpRequest, poolSize), // 使用带缓冲的通道
		logger:     logger,
	}
	engine.initWorkerPool()
	return engine
}

func GetHttpEngine() *HttpEngine {
	if Http == nil {
		Http = NewHttpEngine(200, 10, 3)
	}
	return Http
}

// initWorkerPool 初始化协程池
func (engine *HttpEngine) initWorkerPool() {
	for i := 0; i < engine.poolSize; i++ {
		go engine.worker()
	}
}

// worker 从任务通道中读取任务并执行请求
func (engine *HttpEngine) worker() {
	for req := range engine.tasks {
		engine.wg.Add(1)
		if err := engine.ExecuteRequest(req); err != nil {
			req.Err = err
			engine.logger.Warnf("请求 %d 失败: %v", req.ID, err)
		} else {
			engine.logger.Debugf("请求 %d 成功:", req.ID)
		}
		// 调用回调函数（如果存在）
		if req.Callback != nil {
			req.Callback(req)
		}
		req.Wg.Done()
		engine.wg.Done()
	}
}

// ExecuteRequest 发送HTTP请求
func (engine *HttpEngine) ExecuteRequest(req *HttpRequest) error {
	ctx := context.Background()
	if err := engine.sem.Acquire(ctx, 1); err != nil {
		return err
	}
	defer engine.sem.Release(1)

	var resp *http.Response

	// 构建 HTTP 请求
	httpReq, err := http.NewRequest(req.Method, req.URL, bytes.NewBuffer(req.Body))
	if err != nil {
		return err
	}

	// 设置请求头
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// 发送请求并获取响应
	resp, err = engine.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	response := &HttpResponse{
		raw:     resp,
		code:    resp.StatusCode,
		Headers: resp.Header,
		Body:    body,
	}
	req.Response = response
	req.done = true
	return nil
}

// SubmitRequest 向请求池中添加请求任务
func (engine *HttpEngine) SubmitRequest(req *HttpRequest) {
	req.Wg.Add(1)
	engine.tasks <- req
	req.Wg.Wait()
}

// StopAndWait 等待所有请求完成
func (engine *HttpEngine) StopAndWait() {
	close(engine.tasks) // 关闭任务通道
	engine.wg.Wait()    // 等待所有请求完成
}
