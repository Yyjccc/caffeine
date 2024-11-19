package core

//基础配置

var BasicCfg *BasicConfig

func init() {
	BasicCfg = &BasicConfig{
		ProxyURL: "",
	}
}

type BasicConfig struct {
	ProxyURL string `yaml:"proxy_url"` //http代理
}
