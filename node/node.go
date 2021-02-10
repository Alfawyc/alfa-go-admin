package main

import (
	"go_gin/core/rpc"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//启动节点
	log.Println("rpc serve start 。。。")
	go func() {
		rpc.Start()
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("接收到退出信号，rpc serve quit")
}
