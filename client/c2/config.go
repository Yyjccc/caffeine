package c2

import (
	"caffeine/core"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
)

type C2Yaml struct {
	Name     string     `yaml:"name"`
	Request  C2Request  `yaml:"request"`
	Response C2Response `yaml:"response"`
	Key      CipherKey  `yaml:"key"`
	Basic    C2Basic    `yaml:"basic"`
}

type C2Basic struct {
	Proxy []string `yaml:"proxy"`
}

type ReqCondition struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

// 配置C2请求
type C2Request struct {
	Condition []ReqCondition `yaml:"condition"`
	Method    string         `yaml:"method"`
	//请求加密链
	EncodeChain   string   `yaml:"encode_chain"`
	FrontPadding  string   `yaml:"front_padding"`
	BackPadding   string   `yaml:"back_padding"`
	Headers       []string `yaml:"headers"`
	UserAgentList []string `yaml:"user_agent_list"`
}

// 自定义UnmarshalYAML方法以校验C2Request的字段
func (r *C2Request) UnmarshalYAML(value *yaml.Node) error {
	type rawC2Request C2Request // 定义临时类型以避免递归调用
	var raw rawC2Request
	if err := value.Decode(&raw); err != nil {
		return err
	}
	// 自定义校验
	if raw.Method != "GET" && raw.Method != "POST" && raw.Method != "PUT" && raw.Method != "DELETE" {
		return fmt.Errorf("无效的请求方法:%s", raw.Method)
	}
	if raw.Headers != nil {
		for _, header := range raw.Headers {
			if !core.ValidateExpress(header, ":") {
				return fmt.Errorf("无效header.请以:分割")
			}
		}
	}
	if raw.Condition != nil {
		for _, condition := range raw.Condition {
			if !core.ValidateExpress(condition.Value, "=") {
				return fmt.Errorf("无效表达式.请以=分割")
			}
		}
	}

	*r = C2Request(raw) // 反序列化成功后将值赋给原结构体
	return nil
}

type C2Response struct {
	Headers      []string `yaml:"headers"`
	FrontPadding string   `yaml:"front_padding"`
	BackPadding  string   `yaml:"back_padding"`
	EncodeChain  string   `yaml:"encode_chain"`
	Code         int      `yaml:"code"`
}

type CipherKey struct {
	Xor           string `yaml:"xor"`
	AES           string `yaml:"aes"`
	AESKey        []byte
	XorKey        []byte
	RsaPublicKey  *rsa.PublicKey
	RsaPrivateKey *rsa.PrivateKey
	RsaPublic     string `yaml:"rsa_public"`
	RsaPrivate    string `yaml:"rsa_private"`
}

// 自定义UnmarshalYAML方法以校验CipherKey字段是否为有效的Base64编码和AES密钥有效性
func (c *CipherKey) UnmarshalYAML(value *yaml.Node) error {
	type rawCipherKey CipherKey // 避免递归调用
	var raw rawCipherKey
	if err := value.Decode(&raw); err != nil {
		return err
	}

	if raw.AES != "" {
		c.AES = raw.AES
		// 校验AES密钥
		decoded, err := base64.StdEncoding.DecodeString(raw.AES)
		if err != nil {
			return fmt.Errorf("aes 字段不是有效的 Base64 编码: %v", err)
		} else {
			c.AESKey = decoded
			// 校验 AES 密钥的长度，16、24 或 32 字节为有效长度
			if len(decoded) != 16 && len(decoded) != 24 && len(decoded) != 32 {
				return fmt.Errorf("AES 密钥的长度无效，应为 16、24 或 32 字节")
			}
		}
	}

	if raw.Xor != "" {
		c.Xor = raw.Xor
		decoded, err := base64.StdEncoding.DecodeString(raw.Xor)
		if err != nil {
			return err
		}
		c.XorKey = decoded
	}
	if raw.RsaPublic != "" {
		c.RsaPublic = raw.RsaPublic
		decoded, err := base64.StdEncoding.DecodeString(raw.RsaPublic)
		if err != nil {
			return err
		}
		pemPublicKey := append([]byte("-----BEGIN PUBLIC KEY-----\n"), decoded...)
		pemPublicKey = append(pemPublicKey, []byte("\n-----END PUBLIC KEY-----")...)

		// Step 3: 解码 PEM 格式的公钥
		block, _ := pem.Decode(pemPublicKey)
		if block == nil || block.Type != "PUBLIC KEY" {
			return fmt.Errorf("无效的PEM格式公钥")
		}

		// Step 4: 将解码后的公钥字节数据解析为 rsa.PublicKey 类型
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("解析公钥失败: %v", err)
		}

		// 强制转换为 *rsa.PublicKey 类型
		rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("无法转换为 *rsa.PublicKey 类型")
		}
		c.RsaPublicKey = rsaPublicKey

	}
	if raw.RsaPrivate != "" {
		c.RsaPrivate = raw.RsaPrivate
		decoded, err := base64.StdEncoding.DecodeString(raw.RsaPrivate)
		if err != nil {
			return err
		}
		pemPrivateKey := append([]byte("-----BEGIN RSA PRIVATE KEY-----\n"), decoded...)
		pemPrivateKey = append(pemPrivateKey, []byte("\n-----END RSA PRIVATE KEY-----")...)
		//解码 PEM 格式的私钥
		block, _ := pem.Decode(pemPrivateKey)
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			return fmt.Errorf("无效的PEM格式私钥")
		}
		//  将解码后的私钥字节数据解析为 rsa.PrivateKey 类型
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("解析私钥失败: %v", err)
		}
		c.RsaPrivateKey = privateKey

	}
	// 其他字段可以按类似方式进行校验
	*c = CipherKey(raw)
	return nil
}
