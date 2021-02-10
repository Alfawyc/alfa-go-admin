package system

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go_gin/model"
	"go_gin/model/request"
	"go_gin/model/response"
	"go_gin/service"
	"log"
	"sync"
)

//@Summary 任务列表
//@Tags Task
//@Produce json
//@Param data query request.PageInfo false "页数,每页条数"
//@Router /task/add-task [POST]
func GetTask(ctx *gin.Context) {
	var params request.PageInfo
	_ = ctx.ShouldBindJSON(&params)
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	list, total, err := service.GetTaskList(params)
	if err != nil {
		response.FailWithMessage("获取数据失败"+err.Error(), ctx)
		return
	}
	response.SuccessWithDetail(response.PageResult{
		List:     list,
		Total:    total,
		PageSize: params.PageSize,
		Page:     params.Page,
	}, "success", ctx)
}

//@Summary 添加任务
//@Tags Task
//@Produce json
//@Param data body model.TaskList true "任务信息"
//@Router /task/add-task [POST]
func AddTask(ctx *gin.Context) {
	var task model.TaskList
	_ = ctx.ShouldBindJSON(&task)
	log.Println(task)
	//获取用户
	userId := GetUserId(ctx)
	task.CreatedBy = userId
	//定时任务唯一id
	uuidString := uuid.New().String()
	task.Uuid = uuidString
	res, err := service.AddTaskList(task)
	if err != nil {
		response.FailWithMessage("添加任务失败,"+err.Error(), ctx)
		return
	}
	service.ServiceTask.Add(task)
	response.SuccessWithDetail(gin.H{"task": res}, "success", ctx)
}

//@Summary 停止任务
//@Tags Task
//@Produce json
//@Param id body request.GetById true "任务id"
//@Router /task/stop-task [POST]
func StopTask(ctx *gin.Context) {
	var requestParams request.GetById
	_ = ctx.ShouldBindJSON(&requestParams)
	detail, err := service.TaskDetail(int(requestParams.Id))
	if err != nil {
		response.FailWithMessage("获取任务信息失败"+err.Error(), ctx)
		return
	}
	//暂停任务
	service.ServiceTask.Remove(detail.EntryId)
	response.SuccessWithMessage("success", ctx)
}

//@Summary 重置任务
//@Tags Task
//@Produce json
//@Param id body request.GetById true "任务id"
//@Router /task/recover-task [POST]
func RecoverTask(ctx *gin.Context) {
	var requestParams request.GetById
	_ = ctx.ShouldBindJSON(&requestParams)
	detail, err := service.TaskDetail(int(requestParams.Id))
	if err != nil {
		response.FailWithMessage("获取任务信息失败"+err.Error(), ctx)
		return
	}
	service.ServiceTask.RemoveAndAdd(detail)
	response.SuccessWithMessage("success", ctx)
}

//@Summary 下次运行时间
//@Tags Task
//@Produce json
//@Param id body request.GetById true "任务id"
//@Router /task/next-run [GET]
func NextRun(ctx *gin.Context) {
	var param request.GetById
	_ = ctx.ShouldBindJSON(&param)
	detail, err := service.TaskDetail(int(param.Id))
	if err != nil {
		response.FailWithMessage("获取任务信息失败"+err.Error(), ctx)
		return
	}
	nextTime := service.ServiceTask.NextRunTime(detail)
	response.SuccessWithDetail(gin.H{"next_time": nextTime.Format("2006-01-02 15:04:05")}, "success", ctx)
}

var taskMap sync.Map

//@Summary 终止运行中的任务
//@Tags Task
//@Produce json
//@Param id body request.GetById true "任务id"
//@Router /task/stop-running [POST]
func StopRunning(ctx *gin.Context) {
	var param request.GetById
	_ = ctx.ShouldBindJSON(&param)
	detail, err := service.TaskDetail(int(param.Id))
	if err != nil {
		response.FailWithMessage("获取任务信息失败"+err.Error(), ctx)
		return
	}
	value, ok := taskMap.Load(detail.ID)
	if !ok {
		response.FailWithMessage("任务已执行结束", ctx)
		return
	}
	value.(context.CancelFunc)()
	response.SuccessWithMessage("success", ctx)
	return
}
