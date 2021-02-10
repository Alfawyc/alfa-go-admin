package rpc

import (
	"context"
	"go_gin/core/rpc/proto"
	"go_gin/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

//定义TaskService 服务
type TaskService struct {
}

const (
	//监听地址
	Address string = ":8000"
	//network
	Network string = "tcp"
)

//实现TaskService Run方法
func (s *TaskService) Run(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		}
	}()
	log.Printf("execute cmd start... id: %d , command: %s", req.Id, req.Command)
	output, err := tool.ExecShell(ctx, req.Command)
	resp := new(proto.TaskResponse)
	resp.Output = output
	if err != nil {
		resp.Err = err.Error()
	} else {
		resp.Err = ""
	}
	log.Printf("execute cmd end... id: %d, command: %s, err: %s", req.Id, req.Command, resp.Err)

	return resp, nil
}

//运行serve
func Start() {
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalln("net listen error,", err.Error())
	}
	//输入正式文件和密钥文件为服务端构建TLS凭证
	creds, err := credentials.NewServerTLSFromFile("../core/rpc/key/server.pem", "../core/rpc/key/server.key")
	if err != nil {
		log.Fatalln("[server] generate credentials fail , err", err.Error())
	}
	//新建grpc服务
	grpcServe := grpc.NewServer(grpc.Creds(creds))
	//注册服务
	proto.RegisterTaskServer(grpcServe, &TaskService{})
	//用serve 方法阻塞,等待
	log.Println("start grpcServer with TLS")
	err = grpcServe.Serve(listener)
	if err != nil {
		log.Fatalln("grpcServe err,", err.Error())
	}
}
