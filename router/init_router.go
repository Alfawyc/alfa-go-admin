package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "go_gin/docs"
	"go_gin/middleware"
	"log"
	"net/http"
)

func InitRouter() *gin.Engine {
	//后期修改为自定义中间键
	r := gin.Default()
	r.StaticFS("static/upload/", http.Dir("static/upload"))
	//跨域
	r.Use(middleware.Cors())
	log.Println("user middleware cors")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("register swagger handler")
	//中间件
	//注册系统路由 ,系统路由添加中间件
	r.Use(middleware.Cors())
	InitSystemRouter(r)
	return r
}
