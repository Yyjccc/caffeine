package client

import (
	"caffeine/client/c2"
	"caffeine/client/webshell"
	"caffeine/core"
	"context"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type ClientApp struct {
	ctx     context.Context
	manager *WebShellManger
}

// 导出方法，ui调用
func NewClientApp() *ClientApp {
	return &ClientApp{
		manager: &WebShellManger{
			clients: make(map[int64]*webshell.WebClient),
		},
	}
}

func (a *ClientApp) startup(ctx context.Context) {
	a.ctx = ctx
}

// 获取shell列表
func (a *ClientApp) GetShellList(mode int) []WebShellItem {
	if mode == 0 {
		//本地模式
		return []WebShellItem{
			WebShellItem{
				ID:         0,
				URL:        "http://127.0.0.1/shell.jsp",
				IP:         "127.0.0.1",
				Location:   "湖南",
				Note:       "test",
				CreateTime: time.Now().Format("2006-01-02"),
				UpdateTime: time.Now().Format("2006-01-02"),
				ShellType:  "java",
			},
			WebShellItem{
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
	return []WebShellItem{}
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
	a.manager.AddWebShell(client)
	return client.ID
}

// 测试连接
func (a *ClientApp) TestConnect(id int64) bool {
	client := a.manager.clients[id]
	return client.CheckConnect()
}

// 初始化shell,输出系统信息
func (a *ClientApp) InitShell(id int64) *core.SystemInfo {
	client := a.manager.clients[id]
	client.GetSystemInfo()
	client.LoadDir(client.GetSession().GetCurrentDir())
	return client.GetSession().Info
}

func (a *ClientApp) Exec(id int64, path, cmd string) string {
	client := a.manager.clients[id]
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
