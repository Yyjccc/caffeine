package client

import (
	"caffeine/client/c2"
	"caffeine/client/webshell"
	"caffeine/core"
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/yaml.v3"
)

var (
	App  *ClientApp
	once = sync.Once{}
)

// TerminalManager 终端管理器
type TerminalManager struct {
	terminals map[int64]*webshell.Terminal
}

type ClientApp struct {
	ctx             context.Context
	shellManager    *WebShellManger
	terminalManager *TerminalManager
	taskManger      *webshell.TaskManager
}

// 导出方法，ui调用
func GetClientApp() *ClientApp {
	once.Do(func() {
		appCli := ClientApp{
			shellManager: &WebShellManger{
				clients: make(map[int64]*webshell.WebClient),
			},
			terminalManager: &TerminalManager{
				terminals: make(map[int64]*webshell.Terminal),
			},
			taskManger: webshell.NewTaskManager(),
		}
		App = &appCli
		go App.taskManger.ExecuteAll()
	})
	return App
}

func GetWebClient(sessionID int64) *webshell.WebClient {
	app := GetClientApp()
	return app.shellManager.clients[sessionID]
}

func (a *ClientApp) startup(ctx context.Context) {
	a.ctx = ctx
}

// 获取shell列表
func (a *ClientApp) GetShellList(mode int) []ShellEntry {
	if mode == 0 {
		//本地模式
		return []ShellEntry{
			ShellEntry{
				ID:         0,
				URL:        "http://127.0.0.1/shell.jsp",
				IP:         "127.0.0.1",
				Location:   "湖南",
				Note:       "test",
				CreateTime: time.Now().Format("2006-01-02"),
				UpdateTime: time.Now().Format("2006-01-02"),
				ShellType:  "java",
			},
			ShellEntry{
				ID:         1,
				URL:        "http://127.0.0.1:7878/shell.php",
				IP:         "127.0.0.1",
				Location:   "湖南",
				Note:       "test2",
				CreateTime: time.Now().Format("2006-01-02"),
				UpdateTime: time.Now().Format("2006-01-02"),
				ShellType:  "php",
			},
		}
	}
	return []ShellEntry{}
}

// 进入shell
func (a *ClientApp) GetShellID() int64 {
	//TODO 省略前面数据库操作和新建步骤
	data, err := ioutil.ReadFile("D:\\code\\project\\caffeine\\c2.yaml")
	if err != nil {
		fmt.Errorf("无法读取文件: %v", err)
	}
	//设置代理
	//core.BasicCfg.ProxyURL = "http://127.0.0.1:8083"
	var conf c2.C2Yaml
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Errorf("无法解析 YAML 文件: %v", err)
	}

	target := core.Target{ShellURL: "http://127.0.0.1/shell/server.php"}
	client := webshell.NewWebClient(target, conf)
	a.shellManager.AddWebShell(client)
	return client.ID
}

// 测试连接
func (a *ClientApp) TestConnect(id int64) bool {
	client := a.shellManager.clients[id]
	return client.CheckConnect()
}

// 初始化shell,输出系统信息
func (a *ClientApp) InitShell(id int64) *core.SystemInfo {
	client := a.shellManager.clients[id]
	client.GetSystemInfo()
	//	client.LoadDir(client.GetSession().GetCurrentDir())
	return client.GetSession().Info
}

func (a *ClientApp) Exec(id int64, path, cmd string) string {
	client := a.shellManager.clients[id]
	runCMD := client.RunCMD(path, cmd)
	return runCMD
}

// 获取本地系统状态
func (a *ClientApp) GetLocalSystemMetrics() (*SystemMetric, error) {
	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// 获取内存使用情况
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &SystemMetric{
		CPU:    cpuPercent[0],
		Memory: memInfo.UsedPercent,
		Time:   time.Now().UnixMilli(),
	}, nil
}

// CreateTerminal 创建新终端
func (a *ClientApp) CreateTerminal(shellID int64) (*webshell.TerminalInfo, error) {
	client := a.shellManager.clients[shellID]
	if client == nil {
		return nil, fmt.Errorf("shell client not found: %d", shellID)
	}

	terminal := webshell.NewTerminal(client, client.GetSession().GetCurrentDir())
	a.terminalManager.terminals[terminal.ID] = terminal

	return terminal.GetTerminalInfo(), nil
}

// ExecuteCommand 执行终端命令
func (a *ClientApp) ExecuteCommand(terminalID int64, command string) string {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal == nil {
		return "Terminal not found"
	}

	return terminal.Execute(command)
}

// GetPreviousCommand 获取历史命令
func (a *ClientApp) GetPreviousCommand(terminalID int64) string {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal == nil {
		return ""
	}

	return terminal.GetPreviousCommand()
}

// GetNextCommand 获取下一条历史命令
func (a *ClientApp) GetNextCommand(terminalID int64) string {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal == nil {
		return ""
	}

	return terminal.GetNextCommand()
}

// GetTerminalPrompt 获取终端提示符
func (a *ClientApp) GetTerminalPrompt(terminalID int64) string {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal == nil {
		return "$ "
	}

	return terminal.GetPrompt()
}

// GetTerminalInfo 获取终端信息
func (a *ClientApp) GetTerminalInfo(terminalID int64) *webshell.TerminalInfo {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal == nil {
		return nil
	}

	return terminal.GetTerminalInfo()
}

// CloseTerminal 关闭终端
func (a *ClientApp) CloseTerminal(terminalID int64) error {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal == nil {
		return fmt.Errorf("terminal not found: %d", terminalID)
	}

	delete(a.terminalManager.terminals, terminalID)
	return nil
}

// ListTerminals 列出所有终端
func (a *ClientApp) ListTerminals() []*webshell.TerminalInfo {
	terminals := make([]*webshell.TerminalInfo, 0)
	for _, term := range a.terminalManager.terminals {
		terminals = append(terminals, term.GetTerminalInfo())
	}
	return terminals
}

// GetTerminalWelcomeMessage 获取终端欢迎信息
func (a *ClientApp) GetTerminalWelcomeMessage(terminalID int64) string {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal == nil {
		return "Terminal not found"
	}

	return terminal.GetWelcomeMessage()
}

// SetTerminalEnvironment 设置终端环境变量
func (a *ClientApp) SetTerminalEnvironment(terminalID int64, key, value string) {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal != nil {
		terminal.SetEnvironmentVariable(key, value)
	}
}

// GetTerminalEnvironment 获取终端环境变量
func (a *ClientApp) GetTerminalEnvironment(terminalID int64, key string) string {
	terminal := a.terminalManager.terminals[terminalID]
	if terminal == nil {
		return ""
	}

	return terminal.GetEnvironmentVariable(key)
}

// AddNewShell 添加新的WebShell记录
func (a *ClientApp) AddNewShell(data map[string]interface{}) (int64, error) {
	return a.shellManager.AddNewShell(data)
}

func (a *ClientApp) GetFileSystem(shellID int64) core.FileSystemCache {
	client := a.shellManager.clients[shellID]
	return *client.GetFileSystem()
}

// 加载目录信息
func (a *ClientApp) LoadDirInfo(shellID int64, path string) core.Directory {
	client := a.getWebshellClient(shellID)

	dir := client.LoadDir(path)
	if dir == nil {
		//pass
		return core.Directory{}
	}
	return *dir
}

//// 加载根目录（Windows下所有盘符）
//func (a *ClientApp) LoadFileRoots(shellID int64) core.FileInfo {
//
//}

// 获取操作的实例
func (a *ClientApp) getWebshellClient(shellID int64) *webshell.WebClient {
	return a.shellManager.clients[shellID]
}

// 下载文件,根据文件大小选择方式
func (a *ClientApp) DownloadFile(shellID int64, targetPath string, savePath string) string {
	var task webshell.FileDownloadTask
	//var info *core.FileInfo
	client := GetWebClient(shellID)

	//if info.Size > webshell.DefaultChunkSize {
	//	//分块传输
	//	//pass
	//} else {
	task = webshell.FileDownloadTask{
		Task:           webshell.NewTaskBase(webshell.TaskDownload),
		SessionID:      shellID,
		Client:         client,
		SavePath:       savePath,
		DownloadType:   webshell.SiguredSmall,
		DownloadedSize: 0,
		TargetPath:     targetPath,
	}
	a.taskManger.AddTask(&task)
	//}
	return task.ID
}
