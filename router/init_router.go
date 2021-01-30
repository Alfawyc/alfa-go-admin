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
	//前端页面入口文件
	r.LoadHTMLFiles("dist/index.html")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	}) //检查登陆
	r.StaticFS("static/upload/", http.Dir("static/upload"))
	r.StaticFS("alfa-js/", http.Dir("dist/alfa-js")) //前端项目js文件目录
	//跨域
	r.Use(middleware.Cors())
	log.Println("user middleware cors")
	docUrl := ginSwagger.URL("http://localhost:9191/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, docUrl))
	log.Println("register swagger handler")
	//中间件
	//注册系统路由 ,系统路由添加中间件
	r.Use(middleware.Cors())
	InitSystemRouter(r)
	return r
}
