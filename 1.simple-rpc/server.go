package main

import (
	"context"
	runtime2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-grpc-example/1.simple-rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

type SimpleService struct {
	proto.UnimplementedSimpleServer
}

func (s *SimpleService) Route(ctx context.Context, req *proto.SimpleRequest) (*proto.SimpleResponse, error) {
	res := &proto.SimpleResponse{
		Code:  200,
		Value: "hello " + req.Data,
	}
	return res, nil
}

const (
	// Address 监听地址
	Address string = ":9411"
	// Network 网络通信协议
	Network string = "tcp"
)

func main() {
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	mux := runtime2.NewServeMux()
	proto.RegisterSimpleServer(grpcServer, &SimpleService{})
	err := proto.RegisterSimpleHandlerFromEndpoint(
		context.Background(),
		mux,
		Address,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		panic(err)
	}

	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("grpcServer.Serve err: %v", err)
		}
	}()
	gwServer := &http.Server{
		Addr:    ":9412",
		Handler: mux,
	}
	gwServer.ListenAndServe()

}
