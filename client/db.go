package client

//数据库
//gorm + sqlite

// webshell记录
type WebShellItem struct {
	ID         uint `gorm:"primaryKey"`
	Location   string
	ShellType  string
	IP         string
	CreateTime string
	UpdateTime string
	URL        string
	Note       string
}
