package system

import (
	"github.com/gin-gonic/gin"
	"go_gin/model"
	"go_gin/model/request"
	"go_gin/model/response"
	"go_gin/service"
	"log"
)

//@desc 获取任务列表
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

func AddTask(ctx *gin.Context) {
	var task model.TaskList
	_ = ctx.ShouldBindJSON(&task)
	log.Println(task)
	//获取用户
	userId := GetUserId(ctx)
	task.CreatedBy = userId
	res, err := service.AddTaskList(task)
	if err != nil {
		response.FailWithMessage("添加任务失败,"+err.Error(), ctx)
		return
	}
	response.SuccessWithDetail(gin.H{"task": res}, "success", ctx)
}
