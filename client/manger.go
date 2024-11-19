package client

import "caffeine/client/webshell"

type WebShellManger struct {
	alive   []uint //存放在线shell
	clients map[int64]*webshell.WebClient
}

func (*WebShellManger) name() {

}

func (m *WebShellManger) AddWebShell(client *webshell.WebClient) {
	m.clients[client.ID] = client
}
