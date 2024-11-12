package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

//编译java程序生成字节码

//生产者-消费者模型

// 编译任务结构
type CompileTask struct {
	Code      string // Java代码
	ClassName string // 类名
	Err       error  //错误
	Done      bool   //是否已经完成
	wg        *sync.WaitGroup
	ClassByte []byte //java字节码
}

func NewCompileTask(className, code string) *CompileTask {
	return &CompileTask{
		ClassName: className,
		Code:      code,
		Done:      false,
		wg:        &sync.WaitGroup{},
	}
}

// JavaCompiler 负责调度任务和管理编译
type JavaCompiler struct {
	tempDirPath string            //临时目录路径
	taskQueue   chan *CompileTask // 任务通道
	errors      chan error        // 错误通道
	wg          sync.WaitGroup    // 用于等待所有任务完成
	logger      *logrus.Logger
	javacPath   string // Java 编译器路径
	maxWorkers  int    // 最大并发工作协程数
}

// NewJavaCompiler 初始化编译器，设置并发数量
func NewJavaCompiler(maxWorkers int) *JavaCompiler {
	// 检查是否存在 Java 编译器
	javacPath, err := exec.LookPath("javac")
	if err != nil {
		// 如果系统 PATH 中没有 javac，使用默认 JRE 路径
		javacPath = "./lib/jre/bin/javac"
		if _, err := os.Stat(javacPath); os.IsNotExist(err) {
			logger.Fatal("Java 编译器 (javac) 不存在且没有默认 JRE 配置在 ./lib/jre/")
		}
	}
	compiler := &JavaCompiler{
		taskQueue:   make(chan *CompileTask),
		logger:      logger,
		tempDirPath: "temp",
		errors:      make(chan error),
		maxWorkers:  maxWorkers,
		javacPath:   javacPath, // 存储 javac 路径
	}
	compiler.start()
	return compiler
}

// start 启动协程池来处理编译任务
func (jc *JavaCompiler) start() {
	// 启动多个worker
	err := jc.ensureTempDir()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < jc.maxWorkers; i++ {
		jc.wg.Add(1)
		go jc.worker()
	}
}

// ensureTempDir 确保存在临时目录，如果不存在则创建该目录
func (jc *JavaCompiler) ensureTempDir() error {
	// 检查 temp 目录是否存在
	if _, err := os.Stat(jc.tempDirPath); os.IsNotExist(err) {
		// 如果目录不存在，则创建它
		if err := os.Mkdir(jc.tempDirPath, 0755); err != nil {
			return fmt.Errorf("创建 ./temp 目录失败: %v", err)
		}
	}
	return nil
}

// deleteTempDir 删除指定的临时目录及其内容
func (jc *JavaCompiler) deleteTempDir() error {
	// 检查目录是否存在
	if _, err := os.Stat(jc.tempDirPath); os.IsNotExist(err) {
		return nil
	}
	// 使用 os.RemoveAll 递归删除目录及其内容
	if err := os.RemoveAll(jc.tempDirPath); err != nil {
		return fmt.Errorf("删除目录 %s 失败: %v", jc.tempDirPath, err)
	}
	return nil
}

// AddTask 添加新的编译任务到任务队列
func (jc *JavaCompiler) AddTask(task *CompileTask) {
	jc.taskQueue <- task
}

func (jc *JavaCompiler) compile(task *CompileTask) {
	currentPath, err := os.Getwd()
	if err != nil {
		task.Err = err
		task.Done = true
		return
	}
	// 写入Java文件
	javaFilePath := filepath.Join(currentPath, jc.tempDirPath, task.ClassName+".java")
	// 打印路径信息，用于调试
	jc.logger.Debugf("创建临时文件: %s\n", javaFilePath)
	file, err := os.Create(javaFilePath)
	if err != nil {
		task.Err = fmt.Errorf("创建文件失败: %v", err)
		task.Done = true
		return
	}

	// 将字符串内容写入文件
	_, err = file.WriteString(task.Code)
	defer func() {
		// 关闭文件
		file.Close()
		// 删除临时生成的 .java 文件
		err := os.Remove(javaFilePath)
		if err != nil {
			jc.logger.Errorf("删除 .java 文件失败: %v", err)
		}
	}()
	if err != nil {
		task.Err = fmt.Errorf("写入文件失败: %v", err)
		task.Done = true
		return
	}

	// 编译Java文件
	compileCmd := exec.Command(jc.javacPath, javaFilePath)
	compileCmd.Dir = jc.tempDirPath
	if output, err := compileCmd.CombinedOutput(); err != nil {
		task.Err = fmt.Errorf("编译失败: %v ,输出: %s", err, string(ConvertToUTF8(output)))
		task.Done = true
		return
	}

	// 读取编译后的 .class 文件
	classFilePath := filepath.Join(jc.tempDirPath, task.ClassName+".class")
	classContent, err := os.ReadFile(classFilePath)
	if err != nil {
		task.Err = fmt.Errorf("读取.class文件失败: %v", err)
		task.Done = true
		return
	}
	defer func() {
		err := os.Remove(classFilePath)
		if err != nil {
			jc.logger.Errorf("删除 .class 文件失败: %v", err)
		}
	}()
	task.Done = true
	task.ClassByte = classContent
}

// worker 执行具体的编译任务
func (jc *JavaCompiler) worker() {
	defer func() {
		if r := recover(); r != nil {
			jc.logger.Errorf("Worker encountered a panic: %v", r)
		}
		jc.wg.Done()
	}()

	for task := range jc.taskQueue {
		//task.wg.Add(1)
		jc.compile(task)
		task.wg.Done()
	}
}

// WaitAndClose 等待所有任务完成并关闭通道
func (jc *JavaCompiler) WaitAndClose() {
	close(jc.taskQueue) // 关闭任务队列
	jc.wg.Wait()        // 等待所有协程完成
	err := jc.deleteTempDir()
	if err != nil {
		jc.logger.Error(err)
	}

}

// CompileAsync 异步编译,不会等待编译结果
func (jc *JavaCompiler) CompileAsync(task *CompileTask, callback func(*CompileTask)) {
	task.wg.Add(1)   // 增加 WaitGroup 的计数
	jc.AddTask(task) // 将任务添加到任务队列
	go func() {
		task.wg.Wait() // 等待编译任务完成
		callback(task) // 执行回调函数
	}()
}

// Compile 编译，会等待编译结果
func (jc *JavaCompiler) Compile(task *CompileTask) *CompileTask {
	task.wg.Add(1)
	jc.AddTask(task)
	task.wg.Wait()
	if task.Err != nil {
		jc.logger.Warn(task.Err)
	}
	return task
}
