package core

import (
	"time"
)

type Operate struct {
	OperateType string
	OperateArgs []interface{}
}

type Target struct {
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
	info           SystemInfo
	file           FileSystem
	Environment    map[string]string // 存储环境变量或其他上下文数据
}
