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
	// 获取密码（复制到剪贴板，不直接展示）
	case "GET":
		getPassword(args)
	// 获取密码，并显示明文
	case "GETD":
		getPasswordDirect(args)
	// 设置密码
	case "SET":
		setPassword(args)
	// 展示所有项
	case "ITEMS":
		showItems()
	// 显示项目的描述信息
	case "DESC":
		describeItem(args)
	// 修改密码
	case "CHPWD":
		changePassword(args)
	// 修改描述
	case "CHDESC":
		changeDescription(args)
	// 删除某个密码项
	case "DEL":
		deleteItem(args)
	// 锁定会话（删除会话信息）
	case "LOCK":
		lock()
	// 命令不支持
	default:
		showNotSupportTip(cmd)
	}
}
