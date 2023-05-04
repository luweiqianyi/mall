# rpc demo
以登录服务为例，搭建环境
## 环境准备
1. 在`mall`的同级目录创建`mallGoPath`目录，在`goland`的`File`->`Settings`中配置`GOPATH`为`mallGoPath`目录所在路径，然后执行下列命令准备grpc环境
    ```shell
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    ```
2. 创建目录：`mall/cmd/login/pb`，存放登录相关的rpc相关的文件
3. `pb`目录下创建登录服务的grpc服务文件`login.proto`
    ```proto
    syntax = "proto3";

    package login;
    
    option go_package = "../pb";
    
    message LoginRequest{
    string username=1;
    string password=2;
    }
    
    message LoginResponse{
    int32 code=1;
    string msg=2;
    }
    
    service LoginService{
    rpc Login(LoginRequest)returns(LoginResponse);
    }
    ```
4. `pb`目录下创建`build.sh`文件，增加如下内容
    ```shell
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative \
        login.proto
    ```
    > `require_unimplemented_servers=false`作用: 在LoginServiceServer中不生成mustEmbedUnimplementedLoginServiceServer接口
5. 在`pb`目录所在路径下执行`build.sh`文件，会自动生成以下两个文件
    * `login.pb.go`
    * `login_grpc.pb.go`
6. 生成的上述两个文件在goland中提示`Cannot resolve symbol 'protobuf'`,使用`goland`的`Sync dependency of`机制进行同步即可
7. 在`login`目录下创建`service`目录，创建`login.go`文件，实现文件`login_grpc.pb.go`中的`LoginServiceServer`接口中的`Login`方法，代码如下：
   ```go
   package service
   
   import (
       "context"
       "mall/cmd/login/pb"
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
   ```
   > `LoginServiceServer`接口中的`Login`方法就是上面在`login.proto`中定义，由`build.sh`脚本执行自动生成的，
8. 在`main.go`中创建服务监听程序，完成程序主功能
   ```go
   package main
   
   import (
       "google.golang.org/grpc"
       "log"
       "mall/cmd/login/pb"
       "mall/cmd/login/service"
       "net"
   )
   
   func main() {
       s := grpc.NewServer()
       pb.RegisterLoginServiceServer(s, new(service.LoginServiceServerImpl))
   
       port := ":8081"
       listener, err := net.Listen("tcp", port)
       if err != nil {
           log.Printf("listen port %s failed\n", port)
       }
       s.Serve(listener)
   }
   ```
9. `login`目录下创建`test`目录，创建`login_test.go`文件，增加客户端测试代码，功能是向上一步启动的服务监听程序发起远程rpc调用，代码如下：
   ```go
   package test
   
   import (
       "context"
       "google.golang.org/grpc"
       "google.golang.org/grpc/credentials/insecure"
       "log"
       "mall/cmd/login/pb"
       "testing"
   )
   
   func TestRpcLogin(t *testing.T) {
       targetHost := "localhost:8081"
       conn, err := grpc.Dial(targetHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
       if err != nil {
           log.Printf("Dial %s failed,err=%v\n", targetHost, err)
       }
       defer conn.Close()
   
       c := pb.NewLoginServiceClient(conn)
       rpcReply, err := c.Login(context.Background(), &pb.LoginRequest{
           Username: "leebai",
           Password: "123456",
       })
       if err != nil {
           log.Printf("login failed,err=%v\n", err)
       } else {
           log.Printf("login response: %v\n", rpcReply)
       }
   
       c2 := pb.NewLoginServiceClient(conn)
       rpcReply, err = c2.Login(context.Background(), &pb.LoginRequest{
           Username: "zhangsan",
           Password: "123456",
       })
       if err != nil {
           log.Printf("login failed,err=%v\n", err)
       } else {
           log.Printf("login response: %v\n", rpcReply)
       }
   }
   ```
10. 至此就完成了一个登录请求的远程rpc服务调用