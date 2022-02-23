package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/tidwall/gjson"
)

var (
	g_width   int
	g_height  int
	g_view1s  []string
	g_view1c  []ConsoleColor
	g_view2s  []string
	g_view2c  []ConsoleColor
	g_paths   []string
	g_process []ProcessInfo
	g_rundir  string
	g_exe     []*exec.Cmd
	g_exelog  []*os.File
)

func main() {
	g_width, g_height = terminalWindowSize()
	path, err := os.Executable()
	if err != nil {
		log(logF(), LogLevelError, err.Error())
		return
	}
	var fileName string = filepath.Base(path)
	g_rundir = substrTo(path, 0, len(fileName))
	setupCloseHandler()
	log(logF(), LogLevelInfo, "初始目录 "+g_rundir+" 正在初始化 "+fileName)
	log(logF(), LogLevelInfo, "加载配置文件...")
	LoadConfigFile()
	log(logF(), LogLevelOK, "配置文件加载完成。")
	ProcessChk()
	for {
		err := ProcessList()
		if err != nil {
			break
		}
		ProcessListPrint(g_process)
		printC(g_view1s, g_view1c, 0)
		printC(g_view2s, g_view2c, -1)
		if ProcessWhenClose() {
			ProcessChk()
			println("1111")
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
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
		var exePath string = progVal.String()
		g_paths = append(g_paths, exePath)
		// var exeName string = filepath.Base(exePath)
		// var exePathArr []string = strings.Split(exePath, " ")
		// log(logF(), LogLevelInfo, path)
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
