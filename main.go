package main

import (
	"os"
	"./commands"
	"fmt"
	"syscall"
)

func main()  {
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
