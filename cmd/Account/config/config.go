package config

import (
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"mall/pkg/util"
	"sync"
)

// 以下参数由服务发现机制进行动态配置
const (
	serviceHttpPort     = ":9000"
	gDSN                = "root:123456@tcp(localhost:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"
	TargetAuthRpcServer = "localhost:9002"
)

var gConfig *Config

type Config struct {
	mu               sync.RWMutex
	gServiceHttpPort string
	gDB              *gorm.DB
	gRpcClientConn   *grpc.ClientConn
}

func GetGConfig() *Config {
	return gConfig
}

func InitGConfig() {
	if gConfig == nil {
		gConfig = &Config{}
	}
	db := util.InitMySQLDB(gDSN)
	if db == nil {
		panic("init mysql db connection failed")
		return
	}
	rpcClient, err := util.CreateRpcClient(TargetAuthRpcServer)
	if err != nil {
		panic(fmt.Sprintf("create rpc client connection failed,err:%v", err))
		return
	}
	defer ReleaseGConfig()

	gConfig.mu.Lock()
	defer gConfig.mu.Unlock()
	gConfig.gServiceHttpPort = serviceHttpPort
	gConfig.gDB = db
	gConfig.gRpcClientConn = rpcClient
}

func ReleaseGConfig() {
	if gConfig == nil {
		return
	}

	gConfig.mu.Lock()
	defer gConfig.mu.Unlock()
	if gConfig.gRpcClientConn != nil {
		gConfig.gRpcClientConn.Close()
		gConfig.gRpcClientConn = nil
	}
}

func (c *Config) GetServiceHttpPort() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.gServiceHttpPort
}

func (c *Config) GetDB() *gorm.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.gDB
}

func (c *Config) GetRpcClientConn() *grpc.ClientConn {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.gRpcClientConn
}
