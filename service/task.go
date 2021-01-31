package service

import (
	"context"
	"github.com/robfig/cron/v3"
	"go_gin/common/global"
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
		if item.Status != 1 {
			continue
		}
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
	entryId, _ := ServiceCron.AddFunc(item.Spec, taskFunc)
	//开始执行更新运行状态,entryId
	item.RunningState = 1
	item.EntryId = int(entryId)
	global.Db.Updates(&item)

	log.Println("cron entryId: ", entryId)
}

func (task Task) Remove(id int) {
	entryId := cron.EntryID(id)
	ServiceCron.Remove(entryId)
}

func (task Task) RemoveAndAdd(list model.TaskList) {
	task.Remove(list.EntryId)
	task.Add(list)
}

//jonFunc
func CreateJob(taskModel model.TaskList) cron.FuncJob {
	//todo 改为grpc方式调用
	taskFunc := func() {
		//开始执行任务操作
		taskModel.RunningState = 1
		err := global.Db.Updates(&taskModel).Error
		if err != nil {
			log.Println("state 1 error", err.Error())
		} else {
			log.Println("state 1 success")
		}
		log.Println("执行前的操作")
		//执行任务日志记录
		_ = ExecJob(taskModel)
		time.Sleep(time.Second * 5)
		//执行完成
		log.Println("执行后的操作")
		taskModel.RunningState = 2
		err = global.Db.Updates(&taskModel).Error
		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println("更新成功")
		}
		log.Println("task 执行任务完成 ")
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
	var timeout int
	//超时时间
	if taskModel.Timeout == 0 || taskModel.Timeout > 86400 {
		timeout = 86400
	} else {
		timeout = taskModel.Timeout
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	for i < execTimes {
		output, err := tool.ExecShell(ctx, taskModel.Command)
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