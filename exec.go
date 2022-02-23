package main

import (
	"io/ioutil"
	"os"
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

func backgroundRun(path string, arg ...string) (int, error) {
	// arg = append([]string{"/bin/bash", "-c"}, arg...)
	cmd := exec.Command(path, arg...)
	var fileName string = pathToFileName(path)

	file, err := os.Create(g_rundir + "/" + fileName + ".log")
	if err != nil {
		return -1, err
	}
	cmd.Stdout = file
	cmd.Stderr = file

	// stdout, _ := cmd.StdoutPipe()
	err = cmd.Start()
	go func() {
		cmd.Wait()
		// result, _ := ioutil.ReadAll(stdout) // 讀取輸出結果
		// resdata := string(result)
		// println("==", resdata)
		// stdout.Close()
	}()
	// cmd.Stdout = os.Stdout // cmd.Stdout -> stdout
	// cmd.Stderr = os.Stderr // cmd.Stderr -> stderr
	// 也可以重定向檔案 cmd.Stderr= fd (檔案開啟的描述符即可)
	// stdout, _ := cmd.StdoutPipe() //建立輸出管道
	// defer stdout.Close()
	if err != nil {
		return -1, err
	}
	g_exe = append(g_exe, cmd)
	g_exelog = append(g_exelog, file)

	// 當前執行命令 : cmd.Args
	// 當前執行 PID : cmd.Process.Pid
	return cmd.Process.Pid, nil
}
