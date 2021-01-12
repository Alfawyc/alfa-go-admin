package system

import (
	"github.com/gin-gonic/gin"
	"go_gin/model/request"
	"go_gin/model/response"
	"go_gin/service"
	"log"
)

//@summary: 更新角色api权限
//@Route: /update-casbin
func UpdateCasbin(ctx *gin.Context) {
	var params request.CasbinReceive
	_ = ctx.ShouldBindJSON(&params)
	log.Println(params)
	//todo 数据验证

	err := service.UpdateCasbin(params.AuthorityId, params.CasbinInfos)
	if err != nil {
		response.FailWithMessage("更新casbin失败", ctx)
	} else {
		response.SuccessWithMessage("更新成功", ctx)
	}
}

//@Summary: 获取权限列表
func GetPolicyPathByAuthorityId(ctx *gin.Context) {
	var casbin request.CasbinReceive
	err := ctx.ShouldBindJSON(&casbin)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		ctx.Abort()
		return
	}
	//todo 数据验证
	paths := service.GetAuthorityIdPolicy(casbin.AuthorityId)
	response.SuccessWithDetail(paths, "获取成功", ctx)
}
