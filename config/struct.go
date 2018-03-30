package config

import (
	"encoding/xml"
	"fmt"
	"../utils/fileutil"
)

// 配置数据结构
type PwdKeeperConfig struct {
	XMLName         xml.Name `xml:"PasswordKeeper"`
	CertCfg			CertConfig	`xml:"CertCfg"`	// 证书配置
	UserCfg			UserConfig	`xml:"UserCfg"`	// 用户配置
}

// 证书相关配置
type CertConfig struct{
	PrivateKeyFile  string
	PublicKeyFile	string
}

// 用户配置
type UserConfig struct {
	AppDataDir		string      // 数据存储目录
	BackupDir		string		// 备份目录

}


// 设置备份目录
func (cfg *PwdKeeperConfig) SetBackupDir(backupDir string) {
	// set backup dir only when user entered exists path
	if backupDir == "" {
		return
	}
	exists, err := fileutil.IsPathExists(backupDir)
	if err != nil {
		fmt.Println("Check backup dir failed.")
		return
	}
	if !exists {
		fmt.Printf("Backup dir [%s] not exists!", backupDir)
		return
	}
	cfg.UserCfg.BackupDir = backupDir
}