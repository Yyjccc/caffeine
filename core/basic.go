package core

import (
	"sync"
	"time"
)

var (
	// 全局唯一实例
	instance  *BasicConfig
	basicOnce sync.Once
)

// 基础配置
type BasicConfig struct {
	// 代理设置
	Proxy ProxySettings `yaml:"proxy"`

	// 全局超时设置
	Timeout TimeoutSettings `yaml:"timeout"`

	// 实例锁
	mu sync.RWMutex
}

// ProxySettings 代理配置
type ProxySettings struct {
	// HTTP代理
	HTTPProxy  string `yaml:"http_proxy"`  // HTTP代理地址
	HTTPSProxy string `yaml:"https_proxy"` // HTTPS代理地址

	// SOCKS代理
	SocksProxy string `yaml:"socks_proxy"` // SOCKS代理地址
	SocksVer   int    `yaml:"socks_ver"`   // SOCKS版本(4/5)

	// 认证信息
	Username string `yaml:"username"` // 代理用户名
	Password string `yaml:"password"` // 代理密码

	// 代理规则
	NoProxy  []string `yaml:"no_proxy"`  // 不使用代理的地址
	UseProxy []string `yaml:"use_proxy"` // 强制使用代理的地址

	// 代理模式
	Mode string `yaml:"mode"` // direct/proxy/auto/random

	// 高级设置
	Enabled bool `yaml:"enabled"` // 是否启用代理
	Timeout int  `yaml:"timeout"` // 代理超时(秒)
	Retries int  `yaml:"retries"` // 重试次数

	// 代理池设置
	ProxyPool []string `yaml:"proxy_pool"` // 代理池
	PoolMode  string   `yaml:"pool_mode"`  // round-robin/random/weight
}

// TimeoutSettings 超时配置
type TimeoutSettings struct {
	Dial      int `yaml:"dial"`      // 连接超时
	Read      int `yaml:"read"`      // 读取超时
	Write     int `yaml:"write"`     // 写入超时
	KeepAlive int `yaml:"keepalive"` // 保持连接
}

// 默认配置
var defaultConfig = BasicConfig{
	Proxy: ProxySettings{
		Mode:     "direct",
		Enabled:  true,
		Timeout:  30,
		Retries:  3,
		SocksVer: 5,
		NoProxy:  []string{"localhost", "127.0.0.1", "*.local"},
		UseProxy: []string{},
		PoolMode: "round-robin",

		ProxyPool: []string{"http://127.0.0.1:8083"},
	},
	Timeout: TimeoutSettings{
		Dial:      10,
		Read:      30,
		Write:     30,
		KeepAlive: 60,
	},
}

// GetInstance 获取全局唯一实例
func GetInstance() *BasicConfig {
	basicOnce.Do(func() {
		instance = &defaultConfig
	})
	return instance
}

// Reset 重置为默认配置
func (c *BasicConfig) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Proxy = defaultConfig.Proxy
	c.Timeout = defaultConfig.Timeout
}

// Update 更新配置
func (c *BasicConfig) Update(newConfig BasicConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Proxy = newConfig.Proxy
	c.Timeout = newConfig.Timeout
}

// GetProxyURL 根据协议获取代理地址
func (p *ProxySettings) GetProxyURL(protocol string) string {
	if !p.Enabled {
		return ""
	}

	// 如果配置了代理池，从代理池中选择
	if len(p.ProxyPool) > 0 {
		return p.getProxyFromPool()
	}

	switch protocol {
	case "http":
		return p.HTTPProxy
	case "https":
		return p.HTTPSProxy
	case "socks", "socks4", "socks5":
		return p.SocksProxy
	default:
		return ""
	}
}

// getProxyFromPool 从代理池中获取代理地址
func (p *ProxySettings) getProxyFromPool() string {
	if len(p.ProxyPool) == 0 {
		return ""
	}

	// 根据池模式选择代理
	switch p.PoolMode {
	case "random":
		return p.ProxyPool[time.Now().UnixNano()%int64(len(p.ProxyPool))]
	case "round-robin":
		// 这里需要维护一个计数器，简化版本直接用时间
		return p.ProxyPool[time.Now().Second()%len(p.ProxyPool)]
	default:
		return p.ProxyPool[0]
	}
}

// ShouldUseProxy 判断是否应该使用代理
func (p *ProxySettings) ShouldUseProxy(host string) bool {
	if !p.Enabled {
		return false
	}

	// 检查no_proxy列表
	for _, np := range p.NoProxy {
		if matchHost(host, np) {
			return false
		}
	}

	// 检查use_proxy列表
	for _, up := range p.UseProxy {
		if matchHost(host, up) {
			return true
		}
	}

	// 根据模式决定
	switch p.Mode {
	case "proxy":
		return true
	case "auto":
		// 可以添加自动判断逻辑
		return false
	case "random":
		return time.Now().UnixNano()%2 == 0
	default:
		return false
	}
}

// matchHost checks if a host matches a pattern (supports wildcards)
func matchHost(host, pattern string) bool {
	if pattern[0] == '*' {
		return len(host) >= len(pattern)-1 &&
			host[len(host)-len(pattern)+1:] == pattern[1:]
	}
	return host == pattern
}

// 使用示例：
/*
func main() {
	// 获取全局配置实例
	cfg := GetInstance()

	// 使用默认配置
	proxyURL := cfg.Proxy.GetProxyURL("http")

	// 更新配置
	newConfig := BasicConfig{
		Proxy: ProxySettings{
			Enabled: true,
			HTTPProxy: "http://new-proxy:8080",
			Mode: "proxy",
		},
	}
	cfg.Update(newConfig)

	// 重置为默认配置
	cfg.Reset()
}
*/
