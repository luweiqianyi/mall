package main

import (
	"mall/cmd/Account/config"
	"mall/cmd/Account/dao"
	"mall/cmd/Account/router"
	"mall/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.InitGConfig()
	err := dao.CreateTables()
	if err != nil {
		panic(err)
	}

	go func() {
		route := router.NewWebRouter()
		if route != nil {
			route.StartToServe()
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	exitCode := <-ch
	log.PrintLog("Exit Code:%v", exitCode)
}
