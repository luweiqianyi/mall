package util

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateRpcClient(targetRpcServerHost string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(targetRpcServerHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, err
}

func CloseRpcClient(conn *grpc.ClientConn) {
	if conn != nil {
		conn.Close()
	}
}
