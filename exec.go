package main

import (
	"io/ioutil"
	"os/exec"
	"syscall"
)

func quickRun(name string, arg ...string) (string, int, error) {
	cmd := exec.Command(name, arg...)
	// cmd.Stdout = os.Stdout // cmd.Stdout -> stdout
	// cmd.Stderr = os.Stderr // cmd.Stderr -> stderr
	// 也可以重定向檔案 cmd.Stderr= fd (檔案開啟的描述符即可)

	stdout, _ := cmd.StdoutPipe() //建立輸出管道
	defer stdout.Close()
	err := cmd.Start()
	if err != nil {
		return "", -1, err
	}
	// 當前執行命令 : cmd.Args
	// 當前執行 PID : cmd.Process.Pid
	result, _ := ioutil.ReadAll(stdout) // 讀取輸出結果
	resdata := string(result)
	var res int
	err = cmd.Wait()
	if err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			// 獲取命令執行返回狀態，相當於 shell: echo $?
			res = ex.Sys().(syscall.WaitStatus).ExitStatus()
		}
	}
	return resdata, res, nil
}
