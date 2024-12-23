package core

import "fmt"

// SystemInfo 系统信息结构体
type SystemInfo struct {
	ID int64
	// 基础信息
	FileRoot      string `json:"fileRoot"`      // 根目录
	CurrentDir    string `json:"currentDir"`    // 当前目录
	CurrentUser   string `json:"currentUser"`   // 当前用户
	ProcessArch   int    `json:"processArch"`   // 系统位数
	TempDirectory string `json:"tempDirectory"` // 临时目录
	Hostname      string `json:"hostname"`      // 主机名

	// 网络信息
	IpList         []string                `json:"ipList"`         // IP列表
	NetworkIfaces  map[string]NetIfaceInfo `json:"netIfaces"`      // 网卡信息
	ListeningPorts []string                `json:"listeningPorts"` // 监听端口

	// 系统详情
	Os  OSInfo                 `json:"os"`  // 操作系统信息
	Env map[string]interface{} `json:"env"` // 环境变量

	// 进程信息
	ProcessInfo ProcessInfo `json:"processInfo"` // 进程信息

	// Java环境信息
	Java JavaInfo `json:"java"` // Java信息

	// 重要路径
	WebRoot     string   `json:"webRoot"`     // Web根目录
	ConfigPaths []string `json:"configPaths"` // 配置文件路径
	LogPaths    []string `json:"logPaths"`    // 日志文件路径

	// 安全信息
	SecurityInfo SecurityInfo `json:"security"` // 安全相关信息
}

// ProcessInfo 进程详细信息
type ProcessInfo struct {
	Pid        int      `json:"pid"`        // 进程ID
	ParentPid  int      `json:"ppid"`       // 父进程ID
	StartTime  string   `json:"startTime"`  // 启动时间
	CmdLine    string   `json:"cmdLine"`    // 命令行
	WorkingDir string   `json:"workingDir"` // 工作目录
	Owner      string   `json:"owner"`      // 进程所有者
	Memory     int64    `json:"memory"`     // 内存使用
	OpenFiles  []string `json:"openFiles"`  // 打开的文件
}

// SecurityInfo 安全相关信息
type SecurityInfo struct {
	SELinuxEnabled bool     `json:"selinuxEnabled"` // SELinux状态
	AppArmor       bool     `json:"appArmor"`       // AppArmor状态
	Capabilities   []string `json:"capabilities"`   // 进程权限
	SudoersConfig  string   `json:"sudoersConfig"`  // Sudo配置
	SSHKeys        []string `json:"sshKeys"`        // SSH密钥
	CronJobs       []string `json:"cronJobs"`       // 计划任务
	Services       []string `json:"services"`       // 运行的服务
}

// NetIfaceInfo 网络接口信息
type NetIfaceInfo struct {
	Name        string   `json:"name"`        // 接口名称
	MAC         string   `json:"mac"`         // MAC地址
	IPAddresses []string `json:"ipAddresses"` // IP地址
	MTU         int      `json:"mtu"`         // MTU值
	Status      string   `json:"status"`      // 接口状态
}

// OSInfo 操作系统详细信息
type OSInfo struct {
	Name         string `json:"name"`         // 系统名称
	Version      string `json:"version"`      // 系统版本
	Arch         string `json:"arch"`         // 系统架构
	Kernel       string `json:"kernel"`       // 内核版本
	Distribution string `json:"distribution"` // 发行版信息
	Timezone     string `json:"timezone"`     // 时区
	Language     string `json:"language"`     // 系统语言
	LastBoot     string `json:"lastBoot"`     // 上次启动时间
}

// JavaInfo Java环境信息
type JavaInfo struct {
	RuntimeName    string   `json:"runtimeName"`    // 运行环境名称
	VmVersion      string   `json:"vmVersion"`      // VM版本
	VmName         string   `json:"vmName"`         // VM名称
	Home           string   `json:"home"`           // Java主目录
	LibraryPath    string   `json:"libraryPath"`    // 库路径
	Version        string   `json:"version"`        // Java版本
	ClassPath      string   `json:"classPath"`      // 类路径
	TempDir        string   `json:"tempDir"`        // 临时目录
	ExtDirs        []string `json:"extDirs"`        // 扩展目录
	SecurityPolicy string   `json:"securityPolicy"` // 安全策略
}

func (s *SystemInfo) String() string {
	return fmt.Sprintf(`
根目录: %s
当前用户: %s
当前目录: %s
操作系统信息
	名称:  %s
	版本: %s
	发行版: %s
ip列表: %v
`, s.FileRoot, s.CurrentUser, s.CurrentDir, s.Os.Name, s.Os.Version, s.Os.Arch, s.IpList)
}
