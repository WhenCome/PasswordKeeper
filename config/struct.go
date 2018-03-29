package config

import "encoding/xml"

// 配置数据结构
type PwdKeeperConfig struct {
	XMLName         xml.Name `xml:"PasswordKeeper"`
	CertCfg			CertConfig		// 证书配置
	UserCfg			UserConfig		// 用户配置
}

// 证书相关配置
type CertConfig struct{
	XMLName         xml.Name `xml:"Cert"`
	PrivateKeyFile  string
	PublicKeyFile	string
}

// 用户配置
type UserConfig struct {
	XMLName         xml.Name `xml:"Cert"`
	BackupDir		string		// 备份目录
}
