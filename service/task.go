package service

import (
	"github.com/robfig/cron/v3"
	"go_gin/common/global"
	"go_gin/core/rpc"
	"go_gin/core/rpc/proto"
	"go_gin/model"
	"log"
	"sync"
	"time"
)

var (
	ServiceCron *cron.Cron
	ServiceTask Task
	TaskMap     sync.Map
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

func (task Task) NextRunTime(list model.TaskList) time.Time {
	entry := ServiceCron.Entry(cron.EntryID(list.EntryId))
	return entry.Next
}

//jonFunc
func CreateJob(taskModel model.TaskList) cron.FuncJob {
	//todo 改为grpc方式调用
	taskFunc := func() {
		//开始执行任务操作
		taskLogId := beforeExec(taskModel)
		if taskLogId == 0 {
			return
		}
		result := ExecJob(taskModel)
		afterExec(taskModel, result, taskLogId)
		//执行完成
		log.Println("task done, result: ", result.Result)
	}

	return taskFunc
}

func beforeExec(list model.TaskList) int {
	logId := CreateTaskLog(list)
	if logId == 0 {
		log.Println("写入日志失败")
		return 0
	}
	//更新运行状态
	list.RunningState = TASK_RUNNING
	global.Db.Select("running_state").Where("id = ?", list.ID).Updates(&list)

	return logId
}

func afterExec(list model.TaskList, result TaskResult, logId int) {
	err := UpdateTaskLog(logId, result)
	if err != nil {
		log.Println("任务结束,写入日志失败")
		return
	}
	//更新运行状态
	list.RunningState = TASK_STOP
	global.Db.Select("running_state").Where("id = ? ", list.ID).Updates(&list)
}

func CreateTaskLog(list model.TaskList) int {
	var taskLog model.TaskLog
	taskLog.TaskId = int(list.ID)
	taskLog.StartTime = time.Now()
	err := global.Db.Select("task_id", "start_time").Create(&taskLog).Error
	if err != nil {
		return 0
	}
	return taskLog.Id
}

func UpdateTaskLog(taskLogId int, result TaskResult) error {
	var taskLogModel model.TaskLog
	err := global.Db.Where("id = ?", taskLogId).First(&taskLogModel).Error
	if err != nil {
		return err
	}
	taskLogModel.EndTime = time.Now()
	taskLogModel.Result = result.Result
	taskLogModel.RetryTimes = result.RetryTimes
	return global.Db.Updates(&taskLogModel).Error
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
		output, err := RpcRun(taskModel)
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

func RpcRun(taskModel model.TaskList) (string, error) {
	taskRequest := new(proto.TaskRequest)
	taskRequest.Id = int64(taskModel.ID)
	taskRequest.Command = taskModel.Command
	taskRequest.Timeout = int32(taskModel.Timeout)
	output, err := rpc.Exec(taskRequest)

	return output, err
}
