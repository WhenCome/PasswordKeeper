package commands

import (
	"fmt"
	"../utils/encryptutil"
	"../utils/backuputil"
	"../utils/randutil"
	"../utils/timeutil"
	"../utils/envutil"
	"../utils/fileutil"
	"../config"
	"../db/pwditem"
	"log"
	"os"
	"github.com/atotto/clipboard"
	"strings"
)

// 检查会话，用于确保安全操作
func mustCheckSession() {
	if config.Sess.IsValid() {
		return
	}
	fmt.Println("Session invalid! We need confirm your auth again.")
	// 加载配置
	pwdCfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Load config failed: %s . Have you initialized it before?\n", err)
		return
	}
	fmt.Print("Please enter security code : ")
	securityCode := envutil.ReadLine()
	if !verifySecurityCode(pwdCfg, securityCode) {
		fmt.Println("Security verify failed!")
		os.Exit(-1)
	}
	// 更新session信息
	config.Sess.Revalid()
}

// 验证安全码
func verifySecurityCode(cfg *config.PwdKeeperConfig, code string) bool {
	signRs := encryptutil.Md5(code, cfg.UserCfg.Salt)
	if signRs == cfg.UserCfg.SecurityCode {
		return true
	}
	return false
}

// 强制进行安全验证
func mustVerifySecurity() {
	pwdCfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Load config failed: %s \n. Please try again later.", err)
		os.Exit(-1)
	}
	fmt.Print("Please enter security code : ")
	securityCode := envutil.ReadLine()
	if !verifySecurityCode(pwdCfg, securityCode) {
		fmt.Println("Security verify failed!")
		os.Exit(-1)
	}
}

// 显示帮助信息
func showHelp() {
	outputs := []string{
		"NOTE: before you use password keeper, run init first to make a setup.",
		"Password Keeper command list:",
		"\thelp\n\t\tshow commands help",
		"\tinit\n\t\tinitialize app environment and config, including set security code, security code is required in all of the app time",
		"\tsync\n\t\tsync all app file(data and config) to backup directory, this is a good way to make sure that you can recovery your key data in case the app data was lost, but first, you should set your backup directory when initialize the app run time",
		"\tget [item_key]\n\t\tcopy [item_key]'s password to clipboard, this command won't show plain password directly",
		"\tgetd [item_key]\n\t\tget [item_key]'s password and display it directly",
		"\tset [item_key]\n\t\tset [item_key]'s password",
		"\titems\n\t\tshow all password items, it will just show password item and last update time",
		"\tdesc [item_key]\n\t\tshow description of password item",
		"\tchpwd [item_key]\n\t\tchange password of given password item",
		"\tchdesc [item_key]\n\t\tchange description of giver password item",
		"\tdel [item_key]\n\t\tdelete given password item from database, if you synchronized all the app data, then the delete operation can not be recoveried",
		"\tlock\n\t\tlock the current session(current session will be destroyed), this will make you enter your security code, this will be useful when other people can access your computer",
	}
	for _,output := range outputs {
		fmt.Println(output+"\n")
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

	// 如果已经初始化过，则需要验证用户权限
	if fileutil.IsFileExists(config.InitFlagFile) && pwdCfg.UserCfg.SecurityCode != "" {
		fmt.Print("Please enter security code : ")
		securityCode := envutil.ReadLine()
		if !verifySecurityCode(pwdCfg, securityCode) {
			fmt.Println("Security verify failed!")
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

	// 设置安全密码
	securityPwd := ""
	for securityPwd == "" {
		fmt.Print("Enter security code: ")
		securityPwd := envutil.ReadLine()
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
	backupDir := envutil.ReadLine()
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
	mustCheckSession()
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
	fmt.Println("Sync success.")
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

// 获取密码项
func mustQueryPasswordItem(itemKey string) *pwditem.PwdItem {
	passwdItem, err := pwditem.GetByItem(itemKey)
	if err != nil {
		log.Fatalf("query item failed : %s \n", err)
	}
	if passwdItem == nil {
		log.Fatalf("password item not exists")
	}
	return passwdItem
}

// 获取密码
func getPassword(args []string) {
	mustCheckSession()
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
	err = clipboard.WriteAll(decData)
	if err != nil {
		fmt.Printf("Copy password failed : %s \n", err)
		return
	}
	fmt.Println("Password already copied, you can paste it anywhere you want.")
}

// 获取密码并直接展示明文
func getPasswordDirect(args []string) {
	mustCheckSession()
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
	mustCheckSession()
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	// 1. 输入密码（required）
	var password string
	var description string
	for password == "" {
		fmt.Print("Enter password: ")
		password = envutil.ReadLine()
		if password == "" {
			fmt.Println("Error: password can not be empty, please enter again!")
			continue
		}
		break
	}
	// 2. got item description
	fmt.Print("Enter item description: ")
	description = envutil.ReadLine()
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
	mustCheckSession()
	items, err := pwditem.GetItems()
	if err != nil {
		fmt.Printf("Error while query items : %s \n", err)
		return
	}
	fmt.Println(wrapString("Item", 24),"\t",wrapString("Update Time", 24))
	fmt.Println(wrapString("-------------", 24),"\t",wrapString("----------------", 24))
	for _, item := range items {
		fmt.Println(wrapString(item.Item, 24),"\t",wrapString(item.UpdateTime, 24))
	}
}

// 删除密码项
func deleteItem(args []string) {
	mustCheckSession()
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	// confirm user operation
	fmt.Printf("Do you confirm to remove password for %s? (y/n) : ", itemKey)
	input := envutil.ReadChar()
	if strings.ToLower(input) != "y" {
		fmt.Println("Operation terminated!")
		return
	}
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

// 锁定会话（删除会话信息）
func lock() {
	if config.Sess != nil {
		err := config.Sess.Destroy()
		if err != nil {
			fmt.Printf("Lock current session failed : %s \n", err)
			return
		}
	}
	fmt.Println("Session destroyed!")
}

// 描述项目
func descripeItem(args []string) {
	mustCheckSession()
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	passwdItem, err := pwditem.GetByItem(itemKey)
	if err != nil {
		fmt.Printf("query item failed : %s \n", err)
		return
	}
	fmt.Println(passwdItem.Description)
}

// 修改密码
func changePassword(args []string) {
	mustVerifySecurity()
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	passwdItem := mustQueryPasswordItem(itemKey)
	fmt.Printf("Do you confirm to change password for %s? (y/n) : ", itemKey)
	input := envutil.ReadChar()
	if strings.ToLower(input) != "y" {
		fmt.Println("Operation terminated!")
		return
	}
	fmt.Print("Enter your new password : ")
	pwd1 := envutil.ReadLine()
	fmt.Print("Confirm your new password : ")
	pwd2 := envutil.ReadLine()
	if pwd1 != pwd2 {
		fmt.Println("ERROR : password not match!")
		return
	}
	encPwd, err := encryptutil.EncryptData(pwd1)
	if err != nil {
		log.Fatalf("encrypt data failed : %s \n", err)
	}
	passwdItem.Password = encPwd
	passwdItem.UpdateTime = timeutil.GetCurrentFmtTime()
	_, err = passwdItem.UpdateToDb()
	if err != nil {
		log.Fatalf("update data failed : %s \n", err)
	}
	fmt.Println("Password changed.")
}

// 修改描述信息
func changeDescription(args []string) {
	mustCheckSession()
	if len(args) < 1 {
		fmt.Println("Invalid arguments, use [help] command view commands list and usage.")
		return
	}
	itemKey := args[0]
	passwdItem := mustQueryPasswordItem(itemKey)
	fmt.Print("Enter new description : ")
	description := envutil.ReadLine()
	passwdItem.Description = description
	passwdItem.UpdateTime = timeutil.GetCurrentFmtTime()
	_, err := passwdItem.UpdateToDb()
	if err != nil {
		log.Fatalf("update data failed : %s \n", err)
	}
	fmt.Println("Description changed.")
}

// 执行测试命令
func execTest() {
	// encryptutil.GenRsaKey(2048, config.AppDataDir)
}
