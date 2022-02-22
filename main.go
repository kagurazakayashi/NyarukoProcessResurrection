package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/tidwall/gjson"
)

func main() {
	setupCloseHandler()
	log(logF(), LogLevelInfo, "应用程序保活")
	log(logF(), LogLevelDebug, "加载配置文件...")
	LoadConfigFile()
	terminalWindowSize()
}

// LoadConfigFile 載入啟動初始設定
func LoadConfigFile() {
	f, err := os.OpenFile("config.json", os.O_RDONLY, 0600)
	if err != nil {
		log(logF(), LogLevelError, err.Error())
		os.Exit(-1)
	}
	defer f.Close()
	contentByte, err := ioutil.ReadAll(f)
	if err != nil {
		log(logF(), LogLevelError, err.Error())
		os.Exit(-1)
	}
	var fileData string = string(contentByte)
	if !gjson.Valid(fileData) {
		log(logF(), LogLevelError, "JSON DATA ERR")
		os.Exit(-1)
	}
	var result gjson.Result = gjson.Get(fileData, "prog")
	var prog []gjson.Result = result.Array()
	for _, progVal := range prog {
		var progPath string = progVal.String()
		var path string = filepath.Base(progPath)
		fmt.Println(path)
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
