package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	id    int32
	name  string
	cmd   string
	cpu   int64
	mem   int64
	start int64
}

// ProcessList 獲取程序列表
func ProcessList() ([]ProcessInfo, error) {
	var processInfos []ProcessInfo = []ProcessInfo{}
	pids, err := process.Pids()
	if err != nil || len(pids) == 0 {
		return nil, err
	}
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			log(logF(), LogLevelWarning, "NewProcess "+err.Error())
			continue
		}
		pName, err := p.Name()
		if err != nil {
			log(logF(), LogLevelWarning, "pName "+err.Error())
			pName = ""
			continue
		}
		pCmd, err := p.Cmdline()
		if err != nil {
			log(logF(), LogLevelWarning, "pCmd "+err.Error())
			pCmd = ""
			continue
		}
		pStart, err := p.CreateTime()
		if err != nil {
			log(logF(), LogLevelWarning, "pStart "+err.Error())
			pStart = -1
			// continue
		}
		var pCpu int64 = -1
		var pMem int64 = -1
		// windows 暫無法獲得 CPU 和 MEM
		if runtime.GOOS != "windows" {
			pCpuInfo, err := p.CPUPercent()
			if err != nil {
				log(logF(), LogLevelWarning, "pCpu "+err.Error())
				pCpuInfo = 0
				// continue
			} else {
				pCpu = int64(FloatRound(pCpuInfo, 0))
			}
			pMemInfo, err := p.MemoryInfo()
			if err != nil {
				log(logF(), LogLevelWarning, "pMem "+err.Error())
				// continue
			} else {
				pMem = int64(FloatRound(float64(pMemInfo.RSS)/1024.0, 0))
			}
		}
		pInfo := ProcessInfo{
			id:    pid,
			name:  pName,
			cmd:   pCmd,
			cpu:   pCpu,
			mem:   pMem,
			start: pStart,
		}
		processInfos = append(processInfos, pInfo)
	}
	return processInfos, nil
}

// ProcessListPrint 將獲取到的程序列表輸出
// 前置操作步驟：
// processIds, processNames, err := ProcessList()
// if err != nil { return; }
func ProcessListPrint(processInfos []ProcessInfo) {
	g_view1s = []string{
		"\n PID  |      NAME      | Running Time  | CPU% | MEM(KB) | Command",
		"======+================+===============+======+=========+=",
	}
	g_view1c = []ConsoleColor{ConsoleColorCyan, ConsoleColorCyan}
	g_view1s[0] = tabstr(g_view1s[0], "", g_width, false, true)
	g_view1s[1] = tabstr(g_view1s[1], "=", g_width, false, true)
	for _, pInfo := range processInfos {
		var line []string = []string{}
		line = append(line, tabstr(strconv.FormatInt(int64(pInfo.id), 10), "", 5, true, true))
		line = append(line, tabstr(pInfo.name, "", 14, false, true))
		line = append(line, tabstr(runTime(pInfo.start), "", 13, true, true))
		line = append(line, tabstr(strconv.FormatInt(int64(pInfo.cpu), 10), "", 4, true, true))
		line = append(line, tabstr(strconv.FormatInt(int64(pInfo.mem), 10), "", 7, true, true))
		line = append(line, tabstr(pInfo.cmd, "", g_width-58, false, true))
		var lineStr string = join(line, " | ")
		g_view1s = append(g_view1s, lineStr)
		g_view1c = append(g_view1c, ConsoleColorGreen)
		if len(g_view1s) >= g_height/2 {
			break
		}
	}
}

func runTime(fromTimeStamp int64) string {
	var nowTimeStamp int64 = time.Now().UTC().Unix() * 1000
	var runTimeStamp int64 = nowTimeStamp - fromTimeStamp
	if runTimeStamp > 0 {
		runTimeStamp = runTimeStamp / 1000
	}
	var nDay int64 = runTimeStamp / 3600 / 24
	var nHour int64 = runTimeStamp / 3600
	var nMin int64 = (runTimeStamp % 3600) / 60
	var nSec int64 = (runTimeStamp % 3600) % 60
	return fmt.Sprintf("%d:%02d:%02d:%02d", nDay, nHour, nMin, nSec)
}

// 擷取小數位數
func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}
