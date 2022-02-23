package main

import (
	"fmt"
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
	LogLevelOK      LogLevel = 5
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
	case LogLevelOK:
		return ConsoleColorGreen
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
	// var colorStr string = color.toString()
	// ColorOutput.Colorful.WithFrontColor(colorStr).Println(dateInfo)
	var view2sLen int = len(g_view2s)
	if view2sLen == 0 {
		g_view2s = []string{screenBar()}
		g_view2c = []ConsoleColor{ConsoleColorPurple}
	}
	g_view2s = append(g_view2s, dateInfo)
	g_view2c = append(g_view2c, color)
	if view2sLen > g_height/2-1 {
		g_view2s = append(g_view2s[:0], g_view2s[1:]...)
		g_view2c = append(g_view2c[:0], g_view2c[1:]...)
		g_view2s[0] = screenBar()
		g_view2c[0] = ConsoleColorPurple
	}
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

// substrTo 字串裁剪
// 從 start 開始（包括），到 end 結束（不包括）
// end 為 0 時裁剪到字串末尾
// end 為負數時表示從後至前裁剪多少位
func substrTo(str string, start int, end int) string {
	var strlength int = len(str)
	if strlength == 0 || start > strlength-1 {
		return ""
	}
	var nend int = end
	if end < 0 {
		nend = strlength + end
	}
	if end == 0 || nend > strlength {
		nend = strlength
	}
	if start > nend {
		return ""
	}
	return str[start:nend]
}

// substr 字串裁剪
// 從 start 開始（包括），取 length 长度的字符串
func substr(str string, start int, length int) string {
	return substrTo(str, start, start+length)
}

// tabstr 不足位補齊
// autoSub 超出位裁剪
func tabstr(str string, separator string, toLength int, isRight bool, autoSub bool) string {
	var strlength int = len(str)
	var newStr string = str
	if len(separator) == 0 {
		separator = " "
	}
	if strlength < toLength {
		var addLength int = toLength - strlength
		for i := 0; i < addLength; i++ {
			if isRight {
				newStr = separator + newStr
			} else {
				newStr += separator
			}
		}
	} else if strlength > toLength && autoSub {
		newStr = substr(str, 0, toLength)
	}
	return newStr
}

func printC(strings []string, colors []ConsoleColor, compensate int) {
	var viewLen int = len(strings)
	var height int = g_height / 2
	for i := 0; i < len(strings); i++ {
		ColorOutput.Colorful.WithFrontColor(colors[i].toString()).Println(strings[i])
	}
	if viewLen < height {
		var line string = ""
		for i := 0; i < g_width; i++ {
			line += " "
		}
		for i := 0; i < height-viewLen+compensate; i++ {
			fmt.Println(line)
		}
	}
}

func screenBar() string {
	var line string = ""
	var title string = "▒▒▒ CONSOLE LOG "
	for i := 0; i < g_width-len(title); i++ {
		line += "▒"
	}
	return title + line
}
