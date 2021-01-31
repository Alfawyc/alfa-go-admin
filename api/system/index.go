package system

import (
	"github.com/gin-gonic/gin"
	"go_gin/model/response"
)

//@Summary ping
//@Tags Base
//@Produce json
//@Success 200 {string} json "{"code":200 , "data":"{"data":"Go gin"}" ,"message":"success" }"
//@Router /base/ping [GET]
func HelloWord(ctx *gin.Context) {
	response.SuccessWithDetail(gin.H{"data": "Go Gin"}, "success", ctx)
}

func Check(ctx *gin.Context) {
	userId := GetUserId(ctx)
	if userId == 0 {
		response.FailWithMessage("未登陆", ctx)
		return
	}
	response.SuccessWithMessage("已登陆", ctx)
}
