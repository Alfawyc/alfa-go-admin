package system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go_gin/model/response"
	"go_gin/tool/captcha"
	"net/http"
)

//@Summary 获取验证码
//@Tags Base
//@Produce json
//@Success 200 {string} json "{"code":200 , "data":"" ,"message":"success" }"
//@Router /base/captcha [GET]
func GenerateCaptcha(ctx *gin.Context) {
	id, bs64, err := captcha.DriverDigitGenerate()
	if err != nil {
		response.FailWithMessage("获取验证码失败", ctx)
		return
	}

	response.SuccessWithDetail(gin.H{"data": bs64, "id": id}, "success", ctx)
}

func HtmlCaptcha(ctx *gin.Context) {
	id, bs64, err := captcha.DriverDigitGenerate()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "获取验证码失败"})
	}
	str := fmt.Sprintf("<img src='%s' alt='%s'>", bs64, id)
	ctx.Header("Content-Type", "text/html;charset=utf-8")
	ctx.String(200, str)
}

func VerifyCode(ctx *gin.Context) {
	code := ctx.Param("code")
	id := ctx.Param("id")
	if code == "" || id == "" {
		ctx.JSON(200, gin.H{"message": "缺少验证码和id"})
		return
	}
	store := base64Captcha.DefaultMemStore
	ok := store.Verify(id, code, false)
	var result string
	if ok {
		result = "验证码正确"
	} else {
		result = "验证码错误"
	}
	ctx.JSON(200, gin.H{"code": code, "id": id, "result": result})
}
