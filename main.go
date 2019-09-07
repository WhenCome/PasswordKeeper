package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/whencome/PasswordKeeper/commands"
)

func main() {
	// 设置CPU运行内核数
	cpuNum := runtime.NumCPU()
	if cpuNum > 1 {
		cpuNum--
	}
	runtime.GOMAXPROCS(cpuNum)
	if len(os.Args) < 2 {
		fmt.Println("Invalid command use, use [help] command to view more.")
		syscall.Exit(1)
	}
	command := os.Args[1]
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	commands.ExecCommand(command, args)
}
