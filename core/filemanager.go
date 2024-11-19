package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//文件和目录系统缓存

// Directory 代表一个目录，包含文件和子目录
type Directory struct {
	Name           string       `json:"name"`  // 目录名称
	SubDirectories []*Directory `json:"sub"`   //子目录
	Files          []*FileInfo  `json:"files"` //文件信息
	Path           string       `json:"path"`  //目录绝对路径
	Init           bool
}

// FileInfo 定义单个文件的信息
type FileInfo struct {
	Name         string      `json:"name"`         // 文件名
	Size         int64       `json:"size"`         // 文件大小（字节）
	LastModified time.Time   `json:"lastModified"` // 上次修改时间
	Content      string      // 文件内容
	Permissions  os.FileMode `json:"permissions"` // 文件权限
	FilePath     string      //文件路径
}

// UnmarshalJSON 自定义反序列化逻辑，用于初始化 Directory 和其中的 FilePath
func (d *Directory) UnmarshalJSON(data []byte) error {
	// 定义临时结构用于解析 JSON
	type Alias Directory
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	// 先进行默认的 JSON 反序列化
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 遍历 Directory 中的 Files 字段，初始化每个 FileInfo 的 FilePath
	for _, file := range d.Files {
		d.Path = strings.ReplaceAll(d.Path, "\\", "/")
		file.FilePath = fmt.Sprintf("%s/%s", d.Path, file.Name) // 根据 Directory 的 Path 和 FileName 初始化
	}

	return nil
}

func NewDirectory(path string) *Directory {
	//统一文件分割符
	path = strings.ReplaceAll(path, "\\", "/")
	dirName := filepath.Base(path)
	return &Directory{
		Name: dirName,
		Path: path,
		Init: false,
	}
}

// 构建文件系统管理，方便缓存
type FileSystemCache struct {
	LoadedDirectories map[string]*Directory // 缓存已加载的目录
	Path              string                //当前的绝对路径
	Current           *Directory            // 当前工作目录
	Root              *Directory            // 根目录
}

func NewFileSystem(path string) *FileSystemCache {
	path = strings.ReplaceAll(path, "\\", "/")
	fileSystem := &FileSystemCache{
		LoadedDirectories: make(map[string]*Directory),
		Path:              path,
	}
	current := NewDirectory(path)

	fileSystem.CacheLoadedDir(current)
	fileSystem.Current = fileSystem.LoadedDirectories[current.Path]
	return fileSystem
}

// CacheLoadedDir 缓存/刷新已加载的目录
func (f *FileSystemCache) CacheLoadedDir(dir *Directory) {
	f.LoadedDirectories[dir.Path] = dir
	// 如果当前目录与新增目录的路径一致，更新 Current 的指针
	if f.Current != nil && f.Current.Path == dir.Path {
		f.Current = dir
	}
}

// GetDirectory 根据路径获取Directory对象
func (f *FileSystemCache) GetDirectory(path string) *Directory {
	path = strings.ReplaceAll(path, "\\", "/")
	if path == "." {
		path = f.Current.Path
	}
	if dir, ok := f.LoadedDirectories[path]; ok {
		return dir
	}
	directory := NewDirectory(path)
	f.LoadedDirectories[path] = directory
	return directory
}

// GetFile  从缓存中读取file
func (f *FileSystemCache) GetFile(filePath string) *FileInfo {
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	dirPath := filepath.Dir(filePath)
	directory := f.GetDirectory(dirPath)
	if !directory.Init {
		return nil
	}
	fileName := filepath.Base(filePath)
	for _, file := range directory.Files {
		if fileName == file.Name {
			return file
		}
	}
	return nil
}

// 递归删除目录
func (f *FileSystemCache) RemoveDir(dir *Directory) {
	for _, subDir := range dir.SubDirectories {
		f.RemoveDir(subDir)
		delete(f.LoadedDirectories, subDir.Path)
	}
	delete(f.LoadedDirectories, dir.Path)
}
