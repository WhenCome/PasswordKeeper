package commands

import (
	"fmt"
	"../utils/encryptutil"
	"../utils/backuputil"
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

// 初始化
func initEnv() {
	// 加载配置
	pwdCfg, err := config.LoadConfig()
	if err != nil {
		if err == config.ErrConfigNotExists {
			pwdCfg = &config.PwdKeeperConfig{}
		} else {
			fmt.Printf("Load config failed: %s \n", err)
			return
		}
	}

	// 生成证书信息
	certCfg, err := encryptutil.GenRsaKey(2048, config.AppDataDir)
	if err != nil {
		fmt.Printf("Generate rsa cert failed : %s \n", err)
		return
	}
	pwdCfg.CertCfg = *certCfg

	// 设置数据存储目录
	pwdCfg.UserCfg.AppDataDir = config.AppDataDir
	// 设置备份目录
	fmt.Println("Please enter backup dir:")
	backupDir := ""
	fmt.Scanf("%s", &backupDir)
	pwdCfg.SetBackupDir(backupDir)

	// 写入配置文件
	err = config.SaveConfig(pwdCfg)
	if err != nil {
		fmt.Printf("Error while save config : %s \n", err)
		return
	}
	// 创建初始化标志文件
	config.CreateInitFlag()
}

// 同步配置
func syncConfigs() {
	// 加载配置
	pwdCfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Load config failed: %s \n", err)
		return
	}
	if pwdCfg.UserCfg.BackupDir == "" {
		fmt.Println("Backup dir not set, please use init command initialize first.")
		return
	}
	backuputil.Sync()
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
	// encryptutil.GenRsaKey(2048, config.AppDataDir)
}
