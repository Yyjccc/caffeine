package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

//加密函数

// 支持的加密方式
type CryptoAlgorithm string

const (
	AES    CryptoAlgorithm = "aes"
	RSA    CryptoAlgorithm = "rsa"
	Base64 CryptoAlgorithm = "base64"
	Hex    CryptoAlgorithm = "hex"
	Xor    CryptoAlgorithm = "xor"
	// 可扩展更多算法和编码方式
)

func Base64Encode(src []byte) []byte {
	encoded := base64.StdEncoding.EncodeToString(src)
	return []byte(encoded)
}

func Base64Decode(src []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(src))
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// AESEncode 使用 AES 对输入的明文数据进行加密
// key 为 AES 密钥，长度必须是 16、24 或 32 字节（对应 AES-128、AES-192、AES-256）
// 返回加密后的密文（包含随机 IV）和可能的错误
func AESEncode(plainText, key []byte) ([]byte, error) {
	// 检查密钥长度是否有效
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("AES key length must be 16, 24, or 32 bytes")
	}
	// 创建 AES 加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 创建包含 IV 的加密结果，IV 长度与 AES 块大小相同
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize] // 生成随机 IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	// 创建 CFB 加密模式，并对明文进行加密
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

// AESDecode 使用 AES 对输入的密文数据进行解密
// key 为 AES 密钥，长度必须是 16、24 或 32 字节（对应 AES-128、AES-192、AES-256）
// 返回解密后的明文和可能的错误
func AESDecode(cipherText, key []byte) ([]byte, error) {
	// 检查密钥长度是否有效
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("AES key length must be 16, 24, or 32 bytes")
	}

	// 检查密文长度是否至少包含 IV
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText too short")
	}

	// 创建 AES 解密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 提取 IV（在密文的前 `aes.BlockSize` 个字节）
	iv := cipherText[:aes.BlockSize]
	encryptedText := cipherText[aes.BlockSize:]

	// 创建 CFB 解密流，并对密文进行解密
	stream := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(encryptedText))
	stream.XORKeyStream(plainText, encryptedText)

	return plainText, nil
}

// UnHex 将十六进制编码的字节数组解码成实际的字节数组
func UnHex(src []byte) ([]byte, error) {
	// 检查 src 是否包含偶数个字符，因为每个字节由两个十六进制字符表示
	if len(src)%2 != 0 {
		return nil, errors.New("invalid hex input: length must be even")
	}
	// 创建解码后的字节数组
	decoded := make([]byte, hex.DecodedLen(len(src)))

	// 使用 hex.Decode 将十六进制编码转换为实际的字节
	_, err := hex.Decode(decoded, src)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

// ToHex 将输入字节数组编码为十六进制格式的字节数组
func ToHex(src []byte) []byte {
	// 使用 hex.EncodeToString 将字节数组编码为十六进制字符串
	hexStr := hex.EncodeToString(src)
	// 将十六进制字符串转换为字节数组并返回
	return []byte(hexStr)
}

// xor 对两个字节数组进行按位异或操作
// 如果 src 和 key 长度不同，key 会重复使用，直到与 src 长度相同
func XorCrypto(src, key []byte) []byte {
	result := make([]byte, len(src))

	for i := 0; i < len(src); i++ {
		result[i] = src[i] ^ key[i%len(key)]
	}

	return result
}

// 使用RSA私钥对消息签名
func SignWithRSA(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	// 计算消息的哈希值
	hashed := sha256.Sum256(message)
	// 使用私钥对哈希值签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("签名失败: %v", err)
	}
	return signature, nil
}

// 使用私钥解密数据
func DecryptPrivateRSA(privKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	// 解密：使用私钥解密加密的文本
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
