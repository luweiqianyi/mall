package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"mall/cmd/componentTest/pb"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

type LoginServiceServerImpl struct{}

func (l LoginServiceServerImpl) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	username := request.Username
	password := request.Password
	if username == "leebai" && password == "123456" {
		return &pb.LoginResponse{
			Code: 200,
			Msg:  "success",
		}, nil
	} else {
		return &pb.LoginResponse{
			Code: 201,
			Msg:  "fail",
		}, nil
	}
}

func StartRpcServer() {
	s := grpc.NewServer()
	pb.RegisterLoginServiceServer(s, new(LoginServiceServerImpl))

	port := ":8081"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("listen port %s failed,err=%v\n", port, err)
	}

	s.Serve(listener)
}

// TestRpcLoginService 启动Rpc登录服务
func TestRpcLoginService(t *testing.T) {
	go StartRpcServer()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	exitCode := <-ch
	log.Printf("Exit Code:%v", exitCode)
}

// TestRpcClientLogin Rpc登录客户端，连接启动的Rpc登录服务，发送请求，获取响应
func TestRpcClientLogin(t *testing.T) {
	targetHost := "localhost:8081"
	conn, err := grpc.Dial(targetHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Dial %s failed,err=%v\n", targetHost, err)
	}
	defer conn.Close()

	// 构建一次用户名和密码匹配的rpc请求
	c := pb.NewLoginServiceClient(conn)
	t1 := time.Now()
	rpcReply, err := c.Login(context.Background(), &pb.LoginRequest{
		Username: "leebai",
		Password: "123456",
	})
	log.Printf("spend time:%v", time.Now().Sub(t1).Milliseconds())
	if err != nil {
		log.Printf("login failed,err=%v\n", err)
	} else {
		log.Printf("login response: %v\n", rpcReply)
	}

	//// 构建一次用户名和密码不匹配的rpc请求
	//c2 := pb.NewLoginServiceClient(conn)
	//rpcReply, err = c2.Login(context.Background(), &pb.LoginRequest{
	//	Username: "zhangsan",
	//	Password: "123456",
	//})
	//if err != nil {
	//	log.Printf("login failed,err=%v\n", err)
	//} else {
	//	log.Printf("login response: %v\n", rpcReply)
	//}
}
