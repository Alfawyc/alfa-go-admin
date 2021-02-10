package rpc

import (
	"context"
	"errors"
	"go_gin/core/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"time"
)

var (
	grpcClient proto.TaskClient //grpc 客户端
)

var TaskMap sync.Map

//创建连接
func initClient() {
	//从输入证书为客户端构造TLS凭证
	log.Println("client init")
	cred, err := credentials.NewClientTLSFromFile("core/rpc/key/server.pem",
		"go.ch1451.cn")
	if err != nil {
		log.Fatalln("[client] generate client credentials fail , err ", err.Error())
	}
	//连接服务器
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalln("connect grpcServe fail, ", err.Error())
	}
	//todo 释放连接
	//建立grpc连接
	client := proto.NewTaskClient(conn)
	grpcClient = client
}

//执行任务
func Exec(request *proto.TaskRequest) (string, error) {
	initClient()
	defer func() {
		if err := recover(); err != nil {
			log.Println("rpc client exec error", err)
		}
	}()
	//超时时间判断
	if request.Timeout <= 0 || request.Timeout > 86400 {
		request.Timeout = 86400
	}
	timeout := time.Duration(request.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	TaskMap.Store(int(request.Id), cancel)
	defer TaskMap.Delete(int(request.Id))
	//调用run方法
	resp, err := grpcClient.Run(ctx, request)
	if err != nil {
		statu, ok := status.FromError(err)
		if ok {
			if statu.Code() == codes.DeadlineExceeded {
				log.Println("Run timeout")
			}
		}
		log.Fatalln("grpcClient run error ,", err.Error())
	}
	if resp.Err == "" {
		return resp.Output, nil
	}

	return resp.Output, errors.New(resp.Err)
}
