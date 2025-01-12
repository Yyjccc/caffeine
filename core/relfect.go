package core

import (
	"fmt"
	"reflect"
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

// GetFieldValue 获取结构体字段的值
func GetFieldValue(structPtr interface{}, fieldName string) (interface{}, error) {
	// 获取结构体的反射值
	val := reflect.ValueOf(structPtr)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil, fmt.Errorf("expected a non-nil pointer to a struct, got %v", val.Kind())
	}

	// 获取结构体的值
	val = val.Elem()
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field '%s' not found", fieldName)
	}

	// 返回字段的值
	return field.Interface(), nil
}

// SetFieldValue 设置结构体字段的值
func SetFieldValue(structPtr interface{}, fieldName string, value interface{}) error {
	// 获取结构体的反射值
	val := reflect.ValueOf(structPtr)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("expected a non-nil pointer to a struct, got %v", val.Kind())
	}

	// 获取结构体的值
	val = val.Elem()
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field '%s' not found", fieldName)
	}

	// 确保字段可设置
	if !field.CanSet() {
		return fmt.Errorf("field '%s' cannot be set", fieldName)
	}

	// 将值赋给字段
	fieldValue := reflect.ValueOf(value)
	if fieldValue.Type() != field.Type() {
		return fmt.Errorf("provided value type '%v' doesn't match field type '%v'", fieldValue.Type(), field.Type())
	}

	field.Set(fieldValue)
	return nil
}
