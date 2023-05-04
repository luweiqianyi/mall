package test

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"log"
	"net/http"
	"testing"
)

// 先使用硬编码的方式存储以下值，后期通过配置文件来动态设置以下值
const (
	gHttpsPort = "8087"
	gPemPath   = "./https.pem"
	gKeyPath   = "./https.key"
)

// go get -u github.com/gin-gonic/gin
func TestGin(t *testing.T) {
	eg := gin.Default()
	if eg == nil {
		panic("gin engine initial failed")
	}
	eg.GET("/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Welcome to visit my website",
		})
	})
	eg.Run()
}

// go get -u github.com/unrolled/secure
// 证书的生成下载该网站"https://keymanager.org/"的软件生成即可
func TestGinHttps(t *testing.T) {
	eg := gin.Default()
	if eg == nil {
		panic("gin engine initial failed")
	}
	eg.Use(func(context *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + gHttpsPort,
		})
		err := middleware.Process(context.Writer, context.Request)
		if err != nil {
			log.Printf("handle tls err:%v", err)
			return
		}
		context.Next()
	})
	eg.Run(":"+gHttpsPort, gPemPath, gKeyPath)
}
