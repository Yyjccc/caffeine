package webshell

// 虚拟终端
type Terminal struct {
	//可执行文件路径
	ExecutePath string
	CurrentPath string
}

func NewTerminal(path string) *Terminal {
	return &Terminal{
		CurrentPath: path,
	}
}
