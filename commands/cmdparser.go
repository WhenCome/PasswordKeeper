package commands

import (
	"strings"
)

// 执行命令
func ExecCommand(command string, args []string) {
	cmd := strings.ToUpper(command)
	switch cmd {
	// 显示帮助
	case "HELP":
		showHelp()
	// 同步配置
	case "SYNC":
		syncConfigs()
	// 初始化
	case "INIT":
		initEnv()
	// 获取密码
	case "GET":
		getPassword(args)
	// 设置密码
	case "SET":
		setPassword(args)
	// 测试项目
	case "TEST":
		execTest()
	}
}
