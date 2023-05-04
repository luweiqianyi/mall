package main

import (
	"mall/cmd/Auth/config"
	"mall/cmd/Auth/service"
	"mall/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.InitGConfig()

	go service.StartAuthRpcServer()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	exitCode := <-ch
	log.PrintLog("Exit Code:%v", exitCode)
}
