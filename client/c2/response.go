package c2

import (
	"caffeine/core"
	"container/list"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type ResponseHandler struct {
	config      C2Yaml
	CryptoChain *list.List //解密链，链表结构
}

func NewResponseHandler(config C2Yaml) *ResponseHandler {
	handler := &ResponseHandler{
		config: config,
	}

	handler.parseC2Config(config)
	return handler
}

func (h *ResponseHandler) Handler(session *core.Session, response *core.HttpResponse) ([]byte, error) {
	if response == nil {
		return nil, fmt.Errorf("response is nil")
	}
	body := response.Body
	fmt.Println(string(body))
	//去除填充数据
	if len(h.config.Response.FrontPadding)+len(h.config.Response.BackPadding) > len(body) {
		return nil, fmt.Errorf("Padding length exceeds body length")
	}
	mainData := body[len(h.config.Response.FrontPadding) : len(body)-len(h.config.Response.BackPadding)]
	var err error
	for e := h.CryptoChain.Front(); e != nil; e = e.Next() {
		mainData, err = h.crypto(core.CryptoAlgorithm(e.Value.(string)), mainData)
		if err != nil {
			return nil, err
		}
	}
	return mainData, nil
}

func (h *ResponseHandler) parseC2Config(config C2Yaml) {
	if config.Request.EncodeChain != "" {
		chain := list.New()
		split := strings.Split(config.Response.EncodeChain, "->")
		for _, ciper := range split {
			//往前差入，逆序处理
			chain.PushFront(ciper)
		}
		h.CryptoChain = chain
	}
}

func (h *ResponseHandler) crypto(CryptoAlgorithm core.CryptoAlgorithm, data []byte) ([]byte, error) {
	switch CryptoAlgorithm {
	case core.AES:
		return core.AESDecode(data, h.config.Key.AESKey)
	case core.Base64:
		return base64.StdEncoding.DecodeString(string(data))
	case core.Hex:
		return core.UnHex(data)
	case core.Xor:
		crypto := core.XorCrypto(data, h.config.Key.XorKey)
		return crypto, nil
	case core.RSA:
		//使用私钥解密
		return core.DecryptPrivateRSA(h.config.Key.RsaPrivateKey, data)
	default:
		return nil, errors.New("crypto algorithm not support")
	}
	return data, nil
}
