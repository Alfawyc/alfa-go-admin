package router

import (
	"github.com/gin-gonic/gin"
	"go_gin/api/system"
	"go_gin/middleware"
)

func InitSystemRouter(r *gin.Engine) *gin.RouterGroup {
	//路由组
	g := r.Group("")
	//public route 基础功能 不做鉴权
	publicRoute(g)
	//privateRoute
	privateRoute(g)

	return g
}

func publicRoute(r *gin.RouterGroup) {
	base := r.Group("base")
	{
		base.GET("captcha", system.GenerateCaptcha)
		base.POST("login", system.Login)
		base.GET("ping", system.HelloWord)
	}
	file := r.Group("file")
	{
		file.POST("/upload", system.UploadFile)
	}
}

func privateRoute(r *gin.RouterGroup) {
	privateGroup := r.Group("")
	//中间件
	privateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	InitUserRouter(privateGroup)
	InitCasbinRouter(privateGroup)
	InitAuthRouter(privateGroup)
	InitApiRoute(privateGroup)
}

func InitUserRouter(r *gin.RouterGroup) {
	UserRouter := r.Group("user")
	//UserRouter := r.Group("user").Use(middleware.OperationRecord())
	{
		UserRouter.POST("register", system.Register)
		UserRouter.POST("change-password", system.ChangePassword)
		UserRouter.GET("user-list", system.GetUserList)
		UserRouter.POST("delete-user", system.DeleteUser)
		UserRouter.POST("set-info", system.SetUserInfo)
		UserRouter.POST("user-auth", system.UserAuth)
	}
}

//casbin
func InitCasbinRouter(r *gin.RouterGroup) {
	CasbinRouter := r.Group("casbin")
	{
		CasbinRouter.POST("/update-casbin", system.UpdateCasbin)
		CasbinRouter.POST("/authority-policy", system.GetPolicyPathByAuthorityId)
	}
}

//角色
func InitAuthRouter(r *gin.RouterGroup) {
	AuthRouter := r.Group("auth")
	{
		AuthRouter.POST("/create-auth", system.CreateAuth)
		AuthRouter.POST("/update-auth", system.UpdateAuth)
		AuthRouter.GET("/auth-list", system.AuthList)
		AuthRouter.POST("/delete-auth", system.DeleteAuth)
		AuthRouter.GET("/all-auth", system.AllAuth)
	}
}

func InitApiRoute(r *gin.RouterGroup) {
	ApiRouter := r.Group("api")
	{
		ApiRouter.POST("create-api", system.CreateApi)
		ApiRouter.POST("update-api", system.UpdateApi)
		ApiRouter.POST("delete-api", system.DeleteApi)
		ApiRouter.GET("api-list", system.ApiList)
		ApiRouter.GET("all-api", system.AllApi)
	}
}
