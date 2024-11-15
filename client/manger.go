package client

type WebShellManger struct {
	sessions map[int64]*Session
	alive    []uint //存放在线shell
}
