package commands

import (
	"fmt"
	"../utils/encryptutil"
	"../utils/backuputil"
	"../utils/randutil"
	"../utils/timeutil"
	"../config"
	"../db/pwditem"
	"log"
	"bufio"
	"os"
	"github.com/atotto/clipboard"
)

// 显示帮助信息
func showHelp() {
	outputs := []string{
		"NOTE: before you use password keeper, run init first to make a setup.",
		"Password Keeper command list:",
		"\thelp\n\t\tshow commands help",
		"\tget [item_key]\n\t\tcopy [item_key]'s password to clipboard, this command won't show plain password directly",
		"\tset [item_key]\n\tset [item_key]'s password",
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

	// 用于读取命令行输入
	bio := bufio.NewReader(os.Stdin)

	// 设置安全密码
	securityPwd := ""
	for securityPwd == "" {
		fmt.Print("Enter security code: ")
		line, _, _ := bio.ReadLine()
		securityPwd = string(line)
		// fmt.Scanf("%s", securityPwd)
		if securityPwd == "" {
			fmt.Println("Error: security code can not be empty, please enter again!")
			continue
		}
		pwdCfg.UserCfg.Salt = randutil.GetRandAlphaDigitString(32)
		pwdCfg.UserCfg.SecurityCode = encryptutil.Md5(securityPwd, pwdCfg.UserCfg.Salt)
		break
	}

	// 设置数据存储目录
	pwdCfg.UserCfg.AppDataDir = config.AppDataDir
	// 设置备份目录
	fmt.Println("Please enter backup dir:")
	backupDir := ""
	line, _, _ := bio.ReadLine()
	backupDir = string(line)
	// fmt.Scanf("%s", &backupDir)
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

// 查询密码
func queryPassword(itemKey string) (string, error) {
	passwdItem, err := pwditem.GetByItem(itemKey)
	if err != nil {
		return "", err
	}
	if passwdItem == nil {
		return "",fmt.Errorf("password item [%s] not exists", itemKey)
	}
	decData, err := encryptutil.DecryptData(passwdItem.Password)
	if err != nil {
		return "", err
	}
	return decData, nil
}

// 获取密码
func getPassword(args []string) {
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	/*
	passwdItem, err := pwditem.GetByItem(itemKey)
	if err != nil {
		fmt.Printf("Error while query data from db: %s \n", err)
		return
	}
	if passwdItem == nil {
		fmt.Println("Password item not exists, use set to add it.")
		return
	}
	decData, err := encryptutil.DecryptData(passwdItem.Password)
	if err != nil {
		fmt.Printf("Error wile decrypt data: %s \n", err)
		return
	}
	*/
	decData, err := queryPassword(itemKey)
	if err != nil {
		fmt.Printf("query password failed : %s \n", err)
		return
	}
	err = clipboard.WriteAll(decData)
	if err != nil {
		fmt.Printf("Copy password failed : %s \n", err)
		return
	}
	fmt.Println("Password already copied, you can paste it anywhere you want.")
}

// 获取密码并直接展示明文
func getPasswordDirect(args []string) {
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	decData, err := queryPassword(itemKey)
	if err != nil {
		fmt.Printf("query password failed : %s \n", err)
		return
	}
	fmt.Printf("Your password is : %s \n", decData)
}

// 设置密码
func setPassword(args []string) {
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	// 1. 输入密码（required）
	var password string
	var description string
	bio := bufio.NewReader(os.Stdin)
	// use a loop to get a non-empty password
	for password == "" {
		fmt.Print("Enter password: ")
		line, _, _ := bio.ReadLine()
		password = string(line)
		// fmt.Scanf("%s", &password)
		if password == "" {
			fmt.Println("Error: password can not be empty, please enter again!")
			continue
		}
		break
	}
	// 2. got item description
	fmt.Print("Enter item description: ")
	// fmt.Scanln(&description)
	line, _, _ := bio.ReadLine()
	description = string(line)
	// 3. save password
	passwdItem, err := pwditem.GetByItem(itemKey)
	if err != nil {
		fmt.Printf("Error while query data from db: %s \n", err)
		return
	}
	if passwdItem == nil {
		passwdItem = pwditem.NewPwdItem()
		passwdItem.Item = itemKey
		passwdItem.CreateTime = timeutil.GetCurrentFmtTime()
	}
	passwdItem.UpdateTime = timeutil.GetCurrentFmtTime()
	if description != "" {
		passwdItem.Description = description
	}
	encPwd, err := encryptutil.EncryptData(password)
	if err != nil {
		log.Fatalf("encrypt data failed : %s \n", err)
	}
	passwdItem.Password = encPwd

	// 更新数据
	if passwdItem.Id > 0 {
		_, err = passwdItem.UpdateToDb()
	} else {
		_, err = passwdItem.InsertToDb()
	}
	if err != nil {
		log.Fatalf("Set password for %s failed : %s \n", itemKey, err)
	}
	fmt.Printf("Set password for %s success.\n", itemKey)
}

// 展示所有密码项
func showItems() {
	items, err := pwditem.GetItems()
	if err != nil {
		fmt.Printf("Error while query items : %s \n", err)
		return
	}
	fmt.Println("Item\t\tUpdate Time")
	fmt.Println("-------\t\t---------")
	for _, item := range items {
		fmt.Printf("%s\t\t%s\n", item.Item, item.UpdateTime)
	}
}

// 删除密码项
func deleteItem(args []string) {
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	affectedRows, err := pwditem.DeleteByItem(itemKey)
	if err != nil {
		fmt.Printf("Remove password item failed : %s \n", err)
		return
	}
	if affectedRows == 0 {
		fmt.Printf("Password item not exists. \n")
		return
	}
	fmt.Printf("Password item [%s] already removed. \n", itemKey)
}

// 执行测试命令
func execTest() {
	// encryptutil.GenRsaKey(2048, config.AppDataDir)
}
