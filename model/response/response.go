package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ERROR    = 1001
	SUCCESS  = 0
	JWTError = 999
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Result(code int, data interface{}, msg string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func FailWithDetailed(data interface{}, message string, ctx *gin.Context) {
	Result(ERROR, data, message, ctx)
}

func FailWithMessage(message string, ctx *gin.Context) {
	Result(ERROR, nil, message, ctx)
}

func SuccessWithMessage(message string, ctx *gin.Context) {
	Result(SUCCESS, nil, message, ctx)
}

func SuccessWithDetail(data interface{}, message string, ctx *gin.Context) {
	Result(SUCCESS, data, message, ctx)
}

func JWTFailWithMessage(message string, ctx *gin.Context) {
	Result(JWTError, nil, message, ctx)
}
