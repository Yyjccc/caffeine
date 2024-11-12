package core

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	mu       sync.Mutex
	nextID   int64 = 1  // 起始 ID
	stepSize int64 = 10 // 每次递增的基准步长
)

// GenerateID 生成具有随机性的递增 ID
func GenerateID() int64 {
	mu.Lock()
	defer mu.Unlock()
	// 生成一个在 [1, stepSize-1] 范围内的随机偏移量
	rand.Seed(time.Now().UnixNano())
	offset := int64(rand.Intn(int(stepSize-1)) + 1)
	// 计算随机增量 ID
	id := nextID + offset
	// 更新 nextID，确保生成的 ID 是递增的
	nextID += stepSize
	return id
}

// 自动检测并转换编码
func ConvertToUTF8(input []byte) []byte {
	// 假设首先尝试GBK编码转换
	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Reader := transform.NewReader(bytes.NewReader(input), decoder)
	// 读取并返回转换后的数据
	utf8Output, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		return input
	}
	return utf8Output
}

// 生成一个随机的 (Java 类)名称
func GenerateName() string {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 类名长度
	length := rand.Intn(5) + 5 // 生成 5 到 10 个字符之间的类名

	// Java 类名必须以字母开头，所以确保第一个字符是字母
	var classNameBuilder strings.Builder
	classNameBuilder.WriteByte('A' + byte(rand.Intn(26))) // 随机生成大写字母

	// 随机生成剩下的字符（字母和数字）
	for i := 1; i < length; i++ {
		if rand.Intn(2) == 0 {
			// 随机生成字母
			classNameBuilder.WriteByte('a' + byte(rand.Intn(26)))
		} else {
			// 随机生成数字
			classNameBuilder.WriteByte('0' + byte(rand.Intn(10)))
		}
	}

	return classNameBuilder.String()
}
