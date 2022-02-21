package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	setupCloseHandler()
	log(logF(), LogLevelInfo, "应用程序保活")
}
}

// setupCloseHandler 響應 Ctrl+C
func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log(logF(), LogLevelWarning, "收到中止请求，退出。")
		os.Exit(0)
	}()
}
