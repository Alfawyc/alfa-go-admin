package main

import (
	"context"
	"fmt"
	"go_gin/database"
	"go_gin/router"
	"go_gin/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title Go_admin
// @version 1.0
// @description 这里写描述信息

// @contact.name Alfa
// @contact.email alfa.wang@foxmail.com

// @host http://127.0.0.1:9191
func main() {
	//初始化数据库连接
	database.SetUp()

	r := router.InitRouter()
	srv := &http.Server{
		Addr:    ":9191",
		Handler: r,
	}
	//初始化定时任务
	service.ServiceTask.InitTask()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen：%s \n", err)
		}
	}()
	fmt.Printf("%s Server Run http://127.0.0.1:9191 \r\n", time.Now().Format("2006-01-02 15:04:05"))
	// 优雅的关闭
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server Shutdown")
	//设置超时时间5s
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("Server Shutdown", err)
	}
	log.Println("Server Exiting")
}
