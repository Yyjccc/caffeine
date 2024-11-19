package webshell

// 虚拟终端
type Terminal struct {
	//可执行文件路径
	ExecutePath string
	CurrentPath string
	History     map[string]string
}

func NewTerminal(path string) *Terminal {
	return &Terminal{
		CurrentPath: path,
		History:     make(map[string]string),
	}
}
