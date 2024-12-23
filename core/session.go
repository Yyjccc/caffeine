package core

import (
	"strings"
	"time"
)

// webshell session

type Operate struct {
	OperateType string
	OperateArgs []string
	Time        time.Time
}

type Target struct {
	ID       int64
	ShellURL string
}

// webshell session
type Session struct {
	ID             int64     // 会话唯一标识符
	OperateHistory []Operate // 操作历史
	OutputHistory  []string  // 存储每个命令的执行结果
	StartTime      time.Time // 会话开始时间
	LastActive     time.Time // 上次活跃时间
	Target         Target
	Info           *SystemInfo       //系统信息
	FileSystem     *FileSystemCache  //文件目录缓存
	Environment    map[string]string // 存储环境变量或其他上下文数据
}

// AddOperateHistory  添加操作记录
func (s *Session) AddOperateHistory(funcName string, args []string) {
	operate := Operate{
		OperateType: funcName,
		OperateArgs: args,
		Time:        time.Now(),
	}
	s.OperateHistory = append(s.OperateHistory, operate)
}

// 获取当前目录
func (s *Session) GetCurrentDir() string {
	if s.Info == nil {
		return "./"
	}
	return strings.ReplaceAll(s.Info.CurrentDir, "\\", "/")
}
