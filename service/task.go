package service

import (
	"context"
	"github.com/robfig/cron/v3"
	"go_gin/model"
	"go_gin/tool"
	"log"
	"time"
)

var (
	ServiceCron *cron.Cron
	ServiceTask Task
)

type Task struct {
}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

//初始化定时任务
func (task Task) InitTask() {
	ServiceCron = cron.New(cron.WithSeconds())
	ServiceCron.Start()
	log.Println("开始初始化定时任务")
	//获取所有定时任务
	list, err := GetAllTask()
	if err != nil {
		log.Println("获取定时任务列表错误")
		return
	}
	num := 0
	for _, item := range list {
		task.Add(item)
		num++
	}
	log.Printf("定时任务初始化完成,共添加%d个任务", num)
}

func (task Task) Add(item model.TaskList) {
	if item.Level != 1 {
		log.Println("添加任务失败,子任务无法添加到调度器")
		return
	}
	taskFunc := CreateJob(item)
	if taskFunc == nil {
		log.Println("创建job失败,任务id:", item.ID)
		return
	}
	ServiceCron.AddFunc(item.Spec, taskFunc)
}

func (task Task) Remove(id int) {
	entryId := cron.EntryID(id)
	ServiceCron.Remove(entryId)
}

//jonFunc
func CreateJob(taskModel model.TaskList) cron.FuncJob {
	//todo 改为grpc方式调用
	taskFunc := func() {
		//执行任务日志记录
		result := ExecJob(taskModel)
		log.Println("执行任务结果 , ", result)
	}

	return taskFunc
}

//执行具体任务
func ExecJob(taskModel model.TaskList) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic#task execJob ", err)
		}
	}()
	//默认执行一次
	var execTimes int8 = 1
	if taskModel.RetryTimes > 0 {
		execTimes += taskModel.RetryTimes
	}
	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		output, err := tool.ExecShell(context.Background(), taskModel.Command)
		if err == nil {
			return TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		log.Println("i = ", i)
		if i < execTimes {
			log.Println("任务执行失败,开始重试")
			if taskModel.RetryInterval > 0 {
				time.Sleep(time.Duration(taskModel.RetryInterval) * time.Second)
			} else {
				//默认重试间隔
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}

	return TaskResult{output, err, taskModel.RetryTimes}
}
