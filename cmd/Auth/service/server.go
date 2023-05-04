package service

import (
	"google.golang.org/grpc"
	"mall/cmd/Auth/pb"
	"mall/pkg/log"
	"net"
)

const (
	gRpcListenPort = ":9002"
)

func StartAuthRpcServer() {
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, new(AuthServiceServerImpl))
	listener, err := net.Listen("tcp", gRpcListenPort)
	if err != nil {
		log.PrintLog("listen port %s failed,err=%v\n", gRpcListenPort, err)
	}

	s.Serve(listener)
}
