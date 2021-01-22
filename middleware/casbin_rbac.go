package middleware

import (
	"github.com/gin-gonic/gin"
	"go_gin/model/response"
	"go_gin/service"
	"log"
)

// casbin 中间键
func CasbinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//todo 测试状态
		ctx.Next()
		return
		//获取claims信息
		claims, _ := ctx.Get("claims")
		//类型断言
		waitUser := claims.(*CustomClaims)
		//获取请求url
		obj := ctx.Request.URL.Path //user go to access resource
		//获取请求方法
		act := ctx.Request.Method //the operation that the user performs on the resource
		//获取角色
		sub := waitUser.AuthorityId //用户
		if sub == "" {
			log.Println("暂时未获取到authorityId")
			response.FailWithMessage("未获取到用户角色", ctx)
			return
		}
		e := service.Casbin()
		log.Println(sub, obj, act)
		ok, _ := e.Enforce(sub, obj, act)
		if ok {
			ctx.Next()
		} else {
			response.FailWithDetailed(gin.H{}, "权限不足", ctx)
			ctx.Abort()
			//ctx.Next() //测试阶段
			return
		}
	}
}
