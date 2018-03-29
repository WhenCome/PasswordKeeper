package commands

import (
	"fmt"
	"../utils/encryptutil"
	"../config"
)

// 显示帮助信息
func showHelp() {
	outputs := []string{
		"NOTE: before you use password keeper, run init first to make a setup.",
		"Password Keeper command list:",
		"\thelp\n\t\tshow commands help",
		"\tget [item_key]\n\t\tcopy [item_key]'s password to clipboard, this command won't show plain password directly",
	}
	for _,output := range outputs {
		fmt.Println(output)
	}
}

// 检查工具是否初始化，如果没有初始化，则提示用户
func checkInit() {

}

// 获取密码
func getPassword(args []string) {
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list.")
		return
	}
	itemKey := args[0]
	fmt.Println(itemKey)
}

// 执行测试命令
func execTest() {
	encryptutil.GenRsaKey(2048, config.AppDataDir)
}
