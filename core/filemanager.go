package core

//文件和目录管理系统

// Directory 代表一个目录，包含文件和子目录
type Directory struct {
	Name    string                // 目录名称
	SubDirs map[string]*Directory // 子目录
	Files   map[string]bool       // 目录中的文件
	init    bool
}

func NewDirectory(name string) *Directory {
	return &Directory{
		Name: name,
		init: false,
	}
}

type FileSystem struct {
	Path    string     //当前的绝对路径
	Current *Directory // 当前工作目录
	Root    Directory  // 根目录
}

func NewFileSystem(path string) *FileSystem {
	root := Directory{
		Name:    "/",
		SubDirs: make(map[string]*Directory),
		Files:   make(map[string]bool),
	}
	return &FileSystem{
		Path:    path,
		Current: NewDirectory(path),
		Root:    root,
	}
}

func (f *FileSystem) Cd(path string) {

}
