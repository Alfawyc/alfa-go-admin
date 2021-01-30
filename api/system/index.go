package system

import (
	"github.com/gin-gonic/gin"
	"go_gin/model/response"
	"net/http"
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
	ctx.Redirect(http.StatusMovedPermanently, "/#home/")
}
