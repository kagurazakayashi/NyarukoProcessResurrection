package main

import (
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/phprao/ColorOutput"
)

type ConsoleColor int8

const (
	ConsoleColorBlack  ConsoleColor = 0
	ConsoleColorRed    ConsoleColor = 1
	ConsoleColorGreen  ConsoleColor = 2
	ConsoleColorYellow ConsoleColor = 3
	ConsoleColorBlue   ConsoleColor = 4
	ConsoleColorPurple ConsoleColor = 5
	ConsoleColorCyan   ConsoleColor = 6
	ConsoleColorWhite  ConsoleColor = 7
)

type LogLevel int8

const (
	LogLevelDebug   LogLevel = 0
	LogLevelInfo    LogLevel = 1
	LogLevelWarning LogLevel = 2
	LogLevelError   LogLevel = 3
	LogLevelClash   LogLevel = 4
)

func (p ConsoleColor) toString() string {
	consoleColorString := [8]string{"black", "red", "green", "yellow", "blue", "purple", "cyan", "white"}
	return consoleColorString[p]
}

// LogLevelData 根據日誌記錄級別 lvl 來確定輸出顏色 ConsoleColor
func LogLevelData(lvl LogLevel) ConsoleColor {
	switch lvl {
	case LogLevelDebug:
		return ConsoleColorWhite
	case LogLevelInfo:
		return ConsoleColorCyan
	case LogLevelWarning:
		return ConsoleColorYellow
	case LogLevelError:
		return ConsoleColorRed
	case LogLevelClash:
		return ConsoleColorRed
	default:
		return ConsoleColorCyan
	}
}

// logF 獲取當前函式名稱
func logF() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// log 向終端輸出日誌資訊（根據日誌記錄級別 lvl 來確定輸出顏色）
// 模組名 module 可以從 logF() 獲取
// info 是要顯示的資訊
func log(module string, lvl LogLevel, info string) {
	var timeStr string = time.Now().Format("06-01-02 15:04:05")
	var dateInfo string = "[" + timeStr + "][" + module + "] " + info
	var color ConsoleColor = LogLevelData(lvl)
	var colorStr string = color.toString()
	ColorOutput.Colorful.WithFrontColor(colorStr).Println(dateInfo)
}

// logC 向終端輸出日誌資訊（自定義顏色）
// 提供自定義背景色 ConsoleColor 和自定義前景色 ConsoleColor
func logC(info string, background ConsoleColor, color ConsoleColor) {
	var backgroundStr string = background.toString()
	var colorStr string = color.toString()
	ColorOutput.Colorful.WithFrontColor(colorStr).WithBackColor(backgroundStr).Println(info)
}

// terminalWindowSize 獲取終端顯示區域的縱橫字元容量（寬,高）
func terminalWindowSize() (int, int) {
	var width int = -1
	var height int = -1
	var dataArr []string = []string{}
	if runtime.GOOS == "windows" {
		resdata, res, err := quickRun("powershell.exe", "-noprofile", "-command", "echo", "$host.ui.rawui.WindowSize.Height", "$host.ui.rawui.WindowSize.Width")
		// 不透過 PowerShell : quickRun("mode", "con")
		if err != nil {
			log(logF(), LogLevelWarning, err.Error())
		} else if res != 0 {
			log(logF(), LogLevelWarning, resdata)
		} else {
			dataArr = strings.Split(resdata, "\r\n")
		}
	} else {
		var resStr string = ""
		resdata, res, err := quickRun("tput", "lines")
		if err != nil {
			log(logF(), LogLevelWarning, err.Error())
		} else if res != 0 {
			log(logF(), LogLevelWarning, resdata)
		} else {
			resStr = resdata
			resdata, res, err = quickRun("tput", "cols")
			if err != nil {
				log(logF(), LogLevelWarning, err.Error())
			} else if res != 0 {
				log(logF(), LogLevelWarning, resdata)
			} else {
				resStr = resStr + "\n" + resdata
				dataArr = strings.Split(resStr, "\n")
			}
		}
	}
	for _, line := range dataArr {
		regexpNum := regexp.MustCompile(`\d+`)
		var params []string = regexpNum.FindStringSubmatch(line)
		if len(params) == 0 || len(params[0]) == 0 {
			continue
		}
		lineVal, err := strconv.Atoi(params[0])
		if err != nil {
			continue
		}
		if height == -1 {
			height = lineVal
		} else {
			width = lineVal
			break
		}
	}
	if width < 0 {
		width = 120
	}
	if height < 0 {
		height = 9001
	}
	return width, height
}
