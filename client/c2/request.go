package c2

import (
	"caffeine/core"
	"container/list"
	"encoding/base64"
	"errors"
	"strings"
)

type RequestHandler struct {
	config      C2Yaml
	CryptoChain *list.List //加密链，链表结构
}

func NewRequestHandler(config C2Yaml) *RequestHandler {
	handler := &RequestHandler{
		config: config,
	}

	handler.parseC2Config(config)
	return handler
}

func (h *RequestHandler) parseC2Config(config C2Yaml) {
	if config.Request.EncodeChain != "" {
		chain := list.New()
		split := strings.Split(config.Request.EncodeChain, "->")
		for _, ciper := range split {
			chain.PushBack(ciper)
		}
		h.CryptoChain = chain
	}

}

func (h *RequestHandler) Handler(session core.Session, data []byte) (*core.HttpRequest, error) {
	//加密
	var mainData []byte
	mainData = data
	var err error
	for e := h.CryptoChain.Front(); e != nil; e = e.Next() {
		mainData, err = h.crypto(core.CryptoAlgorithm(e.Value.(string)), mainData)
		if err != nil {
			return nil, err
		}
	}

	req := core.NewHttpRequest()
	req.URL = session.Target.ShellURL
	req.Method = h.config.Request.Method
	reqConfig := h.config.Request
	headers := reqConfig.Headers
	if headers != nil {
		for _, header := range headers {
			split := strings.Split(header, ":")
			if len(split) == 2 {
				req.Headers[split[0]] = split[1]
			}
		}
	}
	//TODO 随机ua头选择

	//预先分配，减少append使用，提高性能
	front := reqConfig.FrontPadding
	back := reqConfig.BackPadding
	body := make([]byte, len(front)+len(mainData)+len(back))
	// 使用 copy 函数将每部分复制到指定位置
	offset := 0
	offset += copy(body[offset:], front)
	offset += copy(body[offset:], mainData)
	copy(body[offset:], back)
	req.Body = body
	return req, nil
}

func (h *RequestHandler) crypto(CryptoAlgorithm core.CryptoAlgorithm, data []byte) ([]byte, error) {
	switch CryptoAlgorithm {
	case core.AES:
		return core.AESEncode(data, h.config.Key.AESKey)
	case core.Base64:
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(encoded, data)
		return encoded, nil
	case core.Hex:
		hex := core.ToHex(data)
		return hex, nil
	case core.Xor:
		crypto := core.XorCrypto(data, h.config.Key.XorKey)
		return crypto, nil
	case core.RSA:
		//使用私钥加密，防止webshell端被截获
		return core.SignWithRSA(h.config.Key.RsaPrivateKey, data)
	default:
		return nil, errors.New("crypto algorithm not support")
	}
	return data, nil
}