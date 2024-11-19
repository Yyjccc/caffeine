package main

import (
	"caffeine/client/c2"
	"caffeine/client/webshell"
	"caffeine/core"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"regexp"
	runtime2 "runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) StartWebShell() *webshell.WebClient {
	data, err := ioutil.ReadFile("../c2.yaml")
	if err != nil {
		fmt.Errorf("无法读取文件: %v", err)
	}
	//设置代理
	core.BasicCfg.ProxyURL = "http://127.0.0.1:8083"
	var conf c2.C2Yaml
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Errorf("无法解析 YAML 文件: %v", err)
	}
	target := core.Target{ShellURL: "http://127.0.0.1/shell/server.php"}
	client := webshell.NewWebClient(target, conf)
	return client
}

func (a *App) Greet() string {
	return "hello"
}

// WailsLogHook 用于捕获日志并将其发送到前端
type WailsLogHook struct {
	ctx context.Context // Wails 上下文，供事件发送使用
}

// NewWailsLogHook 创建一个新的 WailsLogHook
func NewWailsLogHook(ctx context.Context) *WailsLogHook {
	return &WailsLogHook{ctx: ctx}
}

// Levels 定义日志级别，返回所有级别
func (hook *WailsLogHook) Levels() []logrus.Level {
	return logrus.AllLevels // 捕获所有级别的日志
}

// Fire 捕获日志并发送到前端
func (hook *WailsLogHook) Fire(entry *logrus.Entry) error {
	// 获取调用栈信息
	pc, _, line, ok := runtime2.Caller(8)
	var funcName string
	if ok {
		funcName = runtime2.FuncForPC(pc).Name()
	}

	// 创建结构体来存储日志信息
	logData := map[string]interface{}{
		"level":    entry.Level.String(),
		"message":  cleanLogMessage(entry.Message),           // 清理日志内容中的 ANSI 控制字符
		"time":     entry.Time.Format("2006-01-02 15:04:05"), // 时间戳
		"funcName": funcName,                                 // 方法名称
		"line":     line,                                     // 行号
	}

	// 将日志信息转化为 JSON 格式
	logJSON, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	// 使用 Wails 的事件机制将日志消息发送给前端
	runtime.EventsEmit(hook.ctx, "log", string(logJSON))
	return nil
}

// cleanLogMessage 移除 ANSI 控制符的函数
func cleanLogMessage(log string) string {
	// 匹配 ANSI 控制符的正则表达式
	ansiEscapeRegex := regexp.MustCompile(`\033\[[0-9;]*[mK]`)
	// 替换所有匹配到的 ANSI 控制符为空字符串
	return ansiEscapeRegex.ReplaceAllString(log, "")
}
