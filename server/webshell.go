package server

// 加密通信编码器（加密器）
type Encoder interface {
}

type Decoder interface {
}

// 发送到Webshell的数据生成器（不负责加密）
type WebShellServer interface {
	CheckOnline() []byte             //检查是否在线
	GetOsInfo() []byte               //获取系统信息
	RunCmd(path, args string) []byte //运行cmd
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

type CaffeineServer struct {
}

//Caffeine端编码和填充数据
