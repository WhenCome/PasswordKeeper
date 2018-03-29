package config

import (
	"../utils/envutil"
	"../utils/errcode"
	"../utils/fileutil"
	"fmt"
	"os"
)

var AppName string = "PWDKeeper"
var AppDataDir string			// app数据目录
var InitFlagFile string			// 初始化标记
var AppConfigFile string 		// app配置文件

func init() {
	homeDir, err := envutil.GetHomeDir()
	if err != nil {
		fmt.Printf("Got home dir failed: %s\n", err)
		os.Exit(errcode.ERR_GOT_HOME_DIR_FAILED)
	}
	// 初始化数据目录
	AppDataDir = fmt.Sprintf("%s/.config/%s", homeDir, AppName)
	// 检查目录是否存在
	isAppDataDirExists, err := fileutil.IsPathExists(AppDataDir)
	if err != nil {
		fmt.Printf("Got home dir failed: %s\n", err)
		os.Exit(errcode.ERR_STAT_FILE_FAILED)
	}
	if !isAppDataDirExists {
		err = os.MkdirAll(AppDataDir, 0777)
		if err != nil {
			fmt.Printf("create dir [%s] failed: %s\n", AppDataDir, err)
			os.Exit(errcode.ERR_STAT_FILE_FAILED)
		}
	}
	// 至此，创建目录已经成功，开始初始化其他变量（具体文件不创建，需要的时候再创建，只要数据目录存在即可）
	InitFlagFile = fmt.Sprintf("%s/init.dat", AppDataDir)
	AppConfigFile = fmt.Sprintf("%s/app.cfg", AppDataDir)
}