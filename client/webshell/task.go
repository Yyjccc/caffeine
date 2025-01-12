package webshell

import (
	"caffeine/core"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	TaskCreating   = iota // 任务创建中
	TaskInit              // 任务初始化状态完成
	TaskWait              //  任务等待启动
	TaskRunning           // 任务正在执行
	TaskDone              // 任务完成
	TaskFatalError        // 执行中发送致命错误
	TaskTimeout           // 执行超时
	TaskPaused            // 任务暂停
	TaskCancelled         // 任务被取消
)

const (
	TaskDownload = "文件下载"
	TaskUpload   = "文件上传"
)

type Task struct {
	ID         string
	Type       string
	Progress   float64
	Status     int
	CreateTime time.Time
	StartTime  time.Time
	EndTime    time.Time
	TotalTime  time.Duration
	err        error
}

func (t *Task) GetID() string {
	return t.ID
}

func (t *Task) GetType() string {
	return t.Type
}

func (t *Task) SetStartTime(start time.Time) {
	t.StartTime = start
}

func (t *Task) GetError() error {
	return t.err
}

func (t *Task) SetEndTime(end time.Time) {
	t.EndTime = end
	t.TotalTime = end.Sub(t.CreateTime)
}

func (t *Task) UpdateStatus(status int) {
	t.Status = status
}

type RunableTask interface {
	GetID() string
	GetType() string
	GetError() error
	SetStartTime(start time.Time)
	SetEndTime(end time.Time)
	UpdateStatus(status int)
	DoTask() error
	CallBack()
}

func NewTaskBase(Type string) Task {
	return Task{
		ID:         uuid.New().String(),
		Type:       Type,
		Progress:   0,
		CreateTime: time.Now(),
		Status:     TaskCreating,
	}
}

// 文件下载任务
type FileDownloadTask struct {
	Task
	SessionID      int64 //创建任务的 webshell session id
	Client         *WebClient
	FileInfo       *core.FileInfo
	SavePath       string //本地存储路径
	FileName       string //文件名称
	DownloadType   int    //下载类型
	FileSize       int64  //下载文件大小
	DownloadedSize int64  //已经下载的大小
	TargetPath     string //目标机器上面的文件路径
}

const (
	SiguredSmall = iota //加密签名小文件下载
	SiguredBig          //加密签名大文件下载
)

func (task *FileDownloadTask) DoTask() error {

	switch task.DownloadType {
	case SiguredSmall:
		//小文件以文件读取的方式下载
		err := task.Client.DownloadFile(task.TargetPath, task.SavePath)
		if err != nil {
			return err
		}
		task.Progress = 100
	default:
		task.Progress = 100
		return fmt.Errorf("unspoort download type")

	}

	return nil
}

func (task *FileDownloadTask) CallBack() {
}

// 文件上传任务
type FileUploadTask struct {
	Task
}

// 任务管理器
type TaskManager struct {
	tasks      map[string]RunableTask // 存储所有任务
	mu         sync.Mutex             // 任务管理的锁
	submitChan chan RunableTask
	logger     *logrus.Logger
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:      make(map[string]RunableTask),
		submitChan: make(chan RunableTask),
		logger:     core.GetLogger(),
	}
}

// 添加任务
func (tm *TaskManager) AddTask(task RunableTask) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks[task.GetID()] = task
	tm.submitChan <- task
}

// 执行所有任务
func (tm *TaskManager) ExecuteAll() {
	var wg sync.WaitGroup
	for task := range tm.submitChan {
		wg.Add(1)
		go func(task RunableTask) {
			defer wg.Done()
			taskID := task.GetID()
			// 更新任务状态为执行中
			task.UpdateStatus(TaskRunning)
			tm.logger.Infof("启动 %s任务 %s", task.GetType(), task.GetID())
			task.SetStartTime(time.Now())
			err := task.DoTask()
			if err != nil {
				task.UpdateStatus(TaskFatalError)
				tm.logger.Errorf("Error executing task %s: %v", taskID, task)
				return
			}
			task.UpdateStatus(TaskDone)
			task.CallBack()

		}(task)
	}
	wg.Wait()
}
