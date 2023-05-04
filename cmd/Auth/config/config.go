package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"mall/pkg/util"
	"sync"
)

// 以下参数由服务发现机制进行动态配置
const (
	RedisHost     = "localhost:6379"
	RedisPassword = ""
)

var gConfig *Config

type Config struct {
	mu           sync.RWMutex
	gRedisClient *redis.Client
}

func GetGConfig() *Config {
	return gConfig
}

func InitGConfig() {
	if gConfig == nil {
		gConfig = &Config{}
	}
	// TODO 优化：如果redis客户端长期不与Redis服务端进行数据交互，导致客户端连接被服务端关闭的话，可能需要引入心跳机制来对当前
	// Redis连接进行保活
	client := util.NewRedisClient(RedisHost, RedisPassword)
	if client == nil {
		panic(fmt.Sprintf("init redis client connection failed,remote redis addr:%s", RedisHost))
		return
	}
	ReleaseConfig()

	gConfig.mu.Lock()
	defer gConfig.mu.Unlock()
	gConfig.gRedisClient = client
}

func ReleaseConfig() {
	if gConfig == nil {
		return
	}
	if gConfig.gRedisClient != nil {
		gConfig.gRedisClient.Close()
		gConfig.gRedisClient = nil
	}
}

func (c *Config) GetRedisClient() *redis.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.gRedisClient
}

func (c *Config) GetRedisHost() string {
	return RedisHost
}
