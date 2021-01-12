package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_gin/model"
	"go_gin/model/request"
	"go_gin/model/response"
	"go_gin/service"
	"gorm.io/gorm"
)

//@description 创建角色
//@Route /auth/create-auth
func CreateAuth(ctx *gin.Context) {
	var auth model.Auth
	err := ctx.ShouldBind(&auth)
	if err != nil {
		response.FailWithMessage("数据绑定失败,"+err.Error(), ctx)
		return
	}
	activeUserId := GetUserId(ctx)
	auth.CreatedBy = activeUserId
	err, data := service.CreateAuth(auth)
	if err != nil {
		response.FailWithMessage("添加角色失败,"+err.Error(), ctx)
		return
	}
	response.SuccessWithDetail(gin.H{"authority": data}, "添加成功", ctx)
}

//@description 更新角色
//@Route /auth/update-auth
func UpdateAuth(ctx *gin.Context) {
	var auth model.Auth
	_ = ctx.ShouldBind(&auth)
	//todo 数据验证

	err, newAuth := service.UpdateAuthority(auth)
	if err == nil {
		response.SuccessWithDetail(gin.H{"authority": newAuth}, "更新成功", ctx)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.FailWithMessage("角色不存在", ctx)
		return
	}
	response.FailWithMessage("更新角色失败"+err.Error(), ctx)
	return
}

//@description 角色列表
//@Route /auth/auth-list
func AuthList(ctx *gin.Context) {
	var pageInfo request.PageInfo
	_ = ctx.ShouldBind(&pageInfo)
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 10
	}
	err, list, total := service.GetAuthList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取列表失败"+err.Error(), ctx)
		return
	}
	response.SuccessWithDetail(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", ctx)
}

//@description 删除角色
//@Route /auth/delete-auth
func DeleteAuth(ctx *gin.Context) {
	var auth model.Auth
	_ = ctx.ShouldBind(&auth)
	err := service.DeleteAuth(auth)
	if err != nil {
		response.FailWithMessage("删除角色失败,"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("删除成功", ctx)
}

func AllAuth(ctx *gin.Context) {
	var list []model.Auth
	err, list := service.AllAuthList()
	if err != nil {
		response.FailWithMessage("获取角色失败", ctx)
		return
	}

	response.SuccessWithDetail(gin.H{"list": list}, "success", ctx)
}
