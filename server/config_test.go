package server

import (
	"caffeine/client/c2"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestName(t *testing.T) {

	data, err := ioutil.ReadFile("c2.yaml")
	if err != nil {
		fmt.Errorf("无法读取文件: %v", err)
	}

	var conf c2.C2Yaml
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Errorf("无法解析 YAML 文件: %v", err)
	}
	fmt.Println(conf)
}
