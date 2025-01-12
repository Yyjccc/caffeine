package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"github.com/wailsapp/wails/v3/pkg/application"
	"regexp"
	runtime2 "runtime"
)

// App struct
type App struct {
	ctx *application.App
}

// NewApp creates a new App application struct
func NewApp(app *application.App) *App {
	return &App{
		ctx: app,
	}
}

func (a *App) OpenSelectFilePath(fileName string) string {
	options := application.SaveFileDialogOptions{
		CanCreateDirectories:            true,
		ShowHiddenFiles:                 true,
		CanSelectHiddenExtension:        true,
		AllowOtherFileTypes:             true,
		HideExtension:                   false,
		TreatsFilePackagesAsDirectories: true,
		Title:                           "下载文件",
		Message:                         "选择保存路径",
		Directory:                       "",
		Filename:                        fileName,
		ButtonText:                      "保存",
		Window:                          a.ctx.CurrentWindow(),
	}
	saveFileDialog := application.SaveFileDialogWithOptions(&options)
	selection, err := saveFileDialog.PromptForSingleSelection()
	if err != nil {
		return ""
	}
	return selection
}

// WailsLogHook 用于捕获日志并将其发送到前端
type WailsLogHook struct {
	app *application.App // Wails 上下文，供事件发送使用
}

// NewWailsLogHook 创建一个新的 WailsLogHook
func NewWailsLogHook(app *application.App) *WailsLogHook {
	return &WailsLogHook{app: app}
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
	//后端另外一个 协程触发

	// 使用 Wails 的事件机制将日志消息发送给前端
	//EventsEmit(hook.ctx, "log", string(logJSON))
	hook.app.EmitEvent("log", string(logJSON))
	return nil
}

// cleanLogMessage 移除 ANSI 控制符的函数
func cleanLogMessage(log string) string {
	// 匹配 ANSI 控制符的正则表达式
	ansiEscapeRegex := regexp.MustCompile(`\033\[[0-9;]*[mK]`)
	// 替换所有匹配到的 ANSI 控制符为空字符串
	return ansiEscapeRegex.ReplaceAllString(log, "")
}

type TextEdit struct {
}

type UIWindowManager struct{}
