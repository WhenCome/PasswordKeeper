package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/whencome/PasswordKeeper/utils/envutil"
	"github.com/whencome/PasswordKeeper/utils/errcode"
	"github.com/whencome/PasswordKeeper/utils/fileutil"
)

var (
	AppName            = "PasswordKeeper"
	AppDataDir         string           // app数据目录
	InitFlagFile       string           // 初始化标记
	AppConfigFile      string           // app配置文件
	SessionTokenFile   string           // 会话令牌文件
	Sess               *Session         = nil
	AppConfig          *PwdKeeperConfig = nil                                  // 配置信息
	ErrConfigNotExists                  = errors.New("Config File Not Exists") // 错误信息
)

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
	// 数据库相关常量
	DbInitFlagFile = fmt.Sprintf("%s/%s", AppDataDir, DbInitFlag)
	DbFile = fmt.Sprintf("%s/%s", AppDataDir, DbName)
	// 会话相关文件
	SessionTokenFile = fmt.Sprintf("%s/session", AppDataDir)
	// 加载会话信息
	Sess = LoadSession()
	if Sess == nil {
		Sess = NewSession()
	}
}
