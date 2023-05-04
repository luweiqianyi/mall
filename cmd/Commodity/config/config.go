package config

import (
	"gorm.io/gorm"
	"mall/pkg/util"
	"sync"
)

const (
	serviceHttpPort = ":9006"
	gDSN            = "root:123456@tcp(localhost:3306)/mall?charset=utf8mb4&parseTime=True&loc=Local"
)

var gConfig *Config

type Config struct {
	mu               sync.RWMutex
	gServiceHttpPort string
	gDB              *gorm.DB
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
	defer ReleaseGConfig()

	gConfig.mu.Lock()
	defer gConfig.mu.Unlock()
	gConfig.gServiceHttpPort = serviceHttpPort
	gConfig.gDB = db
}

func ReleaseGConfig() {
	if gConfig == nil {
		return
	}
}

func (c *Config) GetDB() *gorm.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.gDB
}

func (c *Config) GetServiceHttpPort() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.gServiceHttpPort
}
