package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

var logger *logrus.Logger

// 定义颜色
const (
	Reset  = "\033[0m"  // 重置
	Red    = "\033[31m" // 红色
	Green  = "\033[32m" // 绿色
	Yellow = "\033[33m" // 黄色
	Blue   = "\033[34m" // 蓝色
)

func init() {
	logger = logrus.New()
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置输出格式
	//logger.SetFormatter(&logrus.TextFormatter{
	//	DisableColors:   false,                 // 启用颜色
	//	FullTimestamp:   true,                  // 显示完整时间戳
	//	TimestampFormat: "2006-01-02 15:04:05", // 自定义时间格式
	//})
	// 设置自定义格式
	logger.SetFormatter(&CustomFormatter{})
}

// CustomFormatter 自定义日志格式
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取调用栈信息
	pc, _, line, ok := runtime.Caller(8)
	var funcName string
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}

	// 使用 path.Base 提取文件名
	//fileName := path.Base(file)

	// 根据日志级别设置颜色
	var levelColor string
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = Blue
	case logrus.InfoLevel:
		levelColor = Green
	case logrus.WarnLevel:
		levelColor = Yellow
	case logrus.ErrorLevel:
		levelColor = Red
	case logrus.FatalLevel:
		levelColor = Red
	case logrus.PanicLevel:
		levelColor = Red
	default:
		levelColor = Reset
	}

	// 格式化日志
	logMessage := fmt.Sprintf("%s[%s]%s %s[%s:%d]%s : %s\n",
		levelColor,
		entry.Level.String(),
		Reset, // 复位颜色
		//fileName,      // 文件名
		Blue,
		funcName, // 函数名
		line,     // 行号
		Reset,
		entry.Message, // 日志消息
	)

	return []byte(logMessage), nil
}
