package core

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestJavaCompiler_AddTask(t *testing.T) {
	compiler := NewJavaCompiler(3)
	code := "package com.yyjccc;" +
		"public class A { public static void main(String[] args) {System.out.println(\"hello,world\");}}"
	task := NewCompileTask("A", code)

	compiler.Compile(task)
	fmt.Println(base64.StdEncoding.EncodeToString(task.ClassByte))
	compiler.WaitAndClose()
}
