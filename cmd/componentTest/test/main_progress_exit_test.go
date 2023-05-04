package test

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestMainProgressExit(t *testing.T) {

	go StartRpcServer()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	exitCode := <-ch
	log.Printf("Exit Code:%v", exitCode)
}
