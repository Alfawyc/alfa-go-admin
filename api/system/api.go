package system

import (
	"github.com/gin-gonic/gin"
	"go_gin/model"
	"go_gin/model/request"
	"go_gin/model/response"
	"go_gin/service"
	"log"
)

//@description 添加api
//@Route /api/create-api
func CreateApi(ctx *gin.Context) {
	var api model.Api
	_ = ctx.ShouldBind(&api)
	log.Println(api)
	err := service.CreateApi(api)
	if err != nil {
		response.FailWithMessage("添加api失败 , "+err.Error(), ctx)
		return
	}
	response.SuccessWithDetail(api, "添加成功", ctx)
}

func DeleteApi(ctx *gin.Context) {
	var params request.GetById
	_ = ctx.ShouldBind(&params)
	err := service.DeleteApi(params.Id)
	if err != nil {
		response.FailWithMessage("删除失败"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("删除成功", ctx)
}

func ApiList(ctx *gin.Context) {
	var params request.PageInfo
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	total, list, err := service.GetApiList(params)
	if err != nil {
		response.FailWithMessage("获取api列表失败", ctx)
		return
	}
	response.SuccessWithDetail(response.PageResult{
		List:     list,
		Total:    total,
		PageSize: params.PageSize,
		Page:     params.Page,
	}, "success", ctx)
}

func AllApi(ctx *gin.Context) {

}
