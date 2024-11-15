package client

import (
	"context"
	"time"
)

type ClientApp struct {
	ctx     context.Context
	manager *WebShellManger
}

// 导出方法，ui调用
func NewClientApp() *ClientApp {
	return &ClientApp{
		manager: &WebShellManger{},
	}
}

func (a *ClientApp) startup(ctx context.Context) {
	a.ctx = ctx
}
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

// 进入webshell
func (a *ClientApp) EnterWebShell(item WebShellItem) (int64, error) {
	return a.manager.EnterShell(item)
}
