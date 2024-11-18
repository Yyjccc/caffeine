package core

import (
	"runtime"
	"strings"
)

// GetCurrentFuncName 获取当前函数名称
func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1) // 获取调用者的信息，1 表示上层调用
	return runtime.FuncForPC(pc).Name()
}

// GetSimpleFuncName 获取当前函数的简单名称（去掉包名）
func GetCallerName() string {
	pc, _, _, _ := runtime.Caller(1)
	fullFuncName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullFuncName, ".")
	return parts[len(parts)-1] // 返回最后的部分，即函数名
}

func GetSimpleFuncName(depth int) string {
	pc, _, _, _ := runtime.Caller(depth)
	fullFuncName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullFuncName, ".")
	return parts[len(parts)-1] // 返回最后的部分，即函数名
}
