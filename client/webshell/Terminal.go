package webshell

import (
	"caffeine/core"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Terminal 虚拟终端,维护执行状态
type Terminal struct {
	ID           int64             // 终端ID
	ExecutePath  string            // shell程序路径
	CurrentPath  string            // 当前工作目录
	Environment  map[string]string // 环境变量
	IsWindows    bool              // 目标系统是否为Windows
	LastExitCode int               // 上一条命令的退出码
	client       *WebClient        // WebShell客户端引用

	// 新增字段
	CommandHistory []string  // 命令历史记录
	HistoryIndex   int       // 当前历史记录索引
	CurrentUser    string    // 当前用户
	StartTime      time.Time // 终端启动时间
}

// TerminalInfo 终端信息结构体(用于前端展示)
type TerminalInfo struct {
	ID          int64  `json:"id"`
	CurrentPath string `json:"currentPath"`
	CurrentUser string `json:"currentUser"`
	IsWindows   bool   `json:"isWindows"`
	ExecutePath string `json:"executePath"`
}

// NewTerminal 创建新的终端实例
func NewTerminal(client *WebClient, path string) *Terminal {
	sysInfo := client.GetSession().Info
	isWindows := strings.Contains(strings.ToLower(sysInfo.Os.Name), "windows")

	t := &Terminal{
		ID:             core.GenerateID(),
		CurrentPath:    path,
		Environment:    make(map[string]string),
		IsWindows:      isWindows,
		LastExitCode:   0,
		client:         client,
		CommandHistory: make([]string, 0),
		HistoryIndex:   -1,
		CurrentUser:    sysInfo.CurrentUser,
		StartTime:      time.Now(),
	}

	t.detectShell()
	return t
}

// GetTerminalInfo 获取终端信息(供前端使用)
func (t *Terminal) GetTerminalInfo() *TerminalInfo {
	return &TerminalInfo{
		ID:          t.ID,
		CurrentPath: t.CurrentPath,
		CurrentUser: t.CurrentUser,
		IsWindows:   t.IsWindows,
		ExecutePath: t.ExecutePath,
	}
}

// Execute 执行命令并返回输出
func (t *Terminal) Execute(cmd string) string {
	// 记录命令历史
	if cmd != "" {
		t.CommandHistory = append(t.CommandHistory, cmd)
		t.HistoryIndex = len(t.CommandHistory)
	}

	// 优先处理内置命令
	if output, handled := t.handleBuiltInCommands(cmd); handled {
		return output
	}

	// 格式化并执行命令
	formattedCmd := t.formatCommand(cmd)
	output := t.client.RunCMD(t.CurrentPath, formattedCmd)

	// 尝试更新退出码
	if t.IsWindows {
		t.LastExitCode = t.getWindowsExitCode()
	} else {
		t.LastExitCode = t.getUnixExitCode()
	}

	return output
}

// GetPreviousCommand 获取历史记录中的上一条命令
func (t *Terminal) GetPreviousCommand() string {
	if len(t.CommandHistory) == 0 || t.HistoryIndex <= 0 {
		return ""
	}
	t.HistoryIndex--
	return t.CommandHistory[t.HistoryIndex]
}

// GetNextCommand 获取历史记录中的下一条命令
func (t *Terminal) GetNextCommand() string {
	if t.HistoryIndex >= len(t.CommandHistory)-1 {
		t.HistoryIndex = len(t.CommandHistory)
		return ""
	}
	t.HistoryIndex++
	return t.CommandHistory[t.HistoryIndex]
}

// getWindowsExitCode 获取Windows下的命令退出码
func (t *Terminal) getWindowsExitCode() int {
	output := t.client.RunCMD(t.CurrentPath, "echo %errorlevel%")
	code := 0
	fmt.Sscanf(output, "%d", &code)
	return code
}

// getUnixExitCode 获取Unix系统下的命令退出码
func (t *Terminal) getUnixExitCode() int {
	output := t.client.RunCMD(t.CurrentPath, "echo $?")
	code := 0
	fmt.Sscanf(output, "%d", &code)
	return code
}

// GetWelcomeMessage 获取欢迎信息
func (t *Terminal) GetWelcomeMessage() string {
	return fmt.Sprintf("Welcome to Virtual Terminal\r\n"+
		"Current User: %s\r\n"+
		"Current Directory: %s\r\n"+
		"Shell: %s\r\n",
		t.CurrentUser,
		t.CurrentPath,
		t.ExecutePath)
}

// detectShell 检测目标系统可用的shell程序
func (t *Terminal) detectShell() {
	if t.IsWindows {
		// 优先检查PowerShell
		if output := t.client.RunCMD(t.CurrentPath, "where powershell.exe"); strings.Contains(output, "powershell.exe") {
			t.ExecutePath = "powershell.exe"
			return
		}
		// 降级使用cmd.exe
		t.ExecutePath = "cmd.exe"
		return
	}

	// 类Unix系统: 优先使用bash,其次是sh
	if output := t.client.RunCMD(t.CurrentPath, "which bash"); output != "" {
		t.ExecutePath = strings.TrimSpace(output)
		return
	}
	if output := t.client.RunCMD(t.CurrentPath, "which sh"); output != "" {
		t.ExecutePath = strings.TrimSpace(output)
		return
	}
	// 默认使用/bin/sh
	t.ExecutePath = "/bin/sh"
}

// handleBuiltInCommands 处理内置命令
func (t *Terminal) handleBuiltInCommands(cmd string) (string, bool) {
	// 处理cd命令
	if strings.HasPrefix(cmd, "cd ") {
		newPath := strings.TrimSpace(strings.TrimPrefix(cmd, "cd "))
		if newPath == "" {
			return "", true
		}

		// 处理相对/绝对路径
		if !filepath.IsAbs(newPath) {
			newPath = filepath.Join(t.CurrentPath, newPath)
		}
		newPath = filepath.Clean(newPath)

		// 在目标系统上验证目录是否存在
		var checkCmd string
		if t.IsWindows {
			checkCmd = fmt.Sprintf("if exist \"%s\\*\" echo Directory_Exists", newPath)
		} else {
			checkCmd = fmt.Sprintf("[ -d \"%s\" ] && echo Directory_Exists", newPath)
		}

		if output := t.client.RunCMD(t.CurrentPath, checkCmd); strings.Contains(output, "Directory_Exists") {
			t.CurrentPath = newPath
			return "", true
		}
		return "目录不存在", true
	}

	// 处理pwd命令
	if cmd == "pwd" {
		return t.CurrentPath, true
	}

	return "", false
}

// formatCommand 根据目标shell类型格式化命令
func (t *Terminal) formatCommand(cmd string) string {
	if t.IsWindows {
		if t.ExecutePath == "powershell.exe" {
			return fmt.Sprintf("%s -Command \"%s\"", t.ExecutePath, cmd)
		}
		return fmt.Sprintf("%s /c %s", t.ExecutePath, cmd)
	}

	// Unix系统命令格式
	return fmt.Sprintf("%s -c '%s'", t.ExecutePath, cmd)
}

// GetPrompt 返回适合的命令提示符
func (t *Terminal) GetPrompt() string {
	if t.IsWindows {
		return fmt.Sprintf("%s> ", t.CurrentPath)
	}
	return fmt.Sprintf("%s$ ", t.CurrentPath)
}

// SetEnvironmentVariable 在目标系统上设置环境变量
func (t *Terminal) SetEnvironmentVariable(key, value string) {
	t.Environment[key] = value
	var cmd string
	if t.IsWindows {
		cmd = fmt.Sprintf("set %s=%s", key, value)
	} else {
		cmd = fmt.Sprintf("export %s=%s", key, value)
	}
	t.Execute(cmd)
}

// GetEnvironmentVariable 从目标系统获取环境变量值
func (t *Terminal) GetEnvironmentVariable(key string) string {
	var cmd string
	if t.IsWindows {
		cmd = fmt.Sprintf("echo %%%s%%", key)
	} else {
		cmd = fmt.Sprintf("echo $%s", key)
	}
	return t.Execute(cmd)
}
