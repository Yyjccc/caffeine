package server

// 加密通信编码器（加密器）
type Encoder interface {
}

type Decoder interface{}

// Webshell 通用接口,所有功能
type WebShell interface {
	CheckOnline() bool    //检查是否在线
	GetOsInfo() string    //获取系统信息
	RunCmd(args []string) //运行cmd
	//下载，上传文件
}

// Webshell 生成器
type WebShellGenerator interface {
}

// webshell 元数据
type ShellMeta struct {
	EncodeChain []string
	DecodeChain []string
}
