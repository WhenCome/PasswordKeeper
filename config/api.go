package config

import (
	"../utils/fileutil"
	"../utils/timeutil"
	"encoding/xml"
)

// 保存配置信息
func SaveConfig(cfg *PwdKeeperConfig) error {
	bytesCfg, err := xml.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = fileutil.WriteFile(AppConfigFile, string(bytesCfg))
	return err
}

// 加载配置信息
func LoadConfig() (*PwdKeeperConfig, error) {
	// 检查配置文件是否存在
	if !fileutil.IsFileExists(AppConfigFile) {
		return nil,ErrConfigNotExists
	}
	// 读取文件内容
	cfgContent,err := fileutil.GetContents(AppConfigFile)
	if err != nil {
		return nil,err
	}
	pwdCfg := &PwdKeeperConfig{}
	err = xml.Unmarshal(cfgContent, pwdCfg)
	if err != nil {
		return nil,err
	}
	return pwdCfg,nil
}
/*
func LoadConfig() (*PwdKeeperConfig, error) {
	cfg, err := LoadConfig()
	if err != nil {
		if err == ErrConfigNotExists {
			return &PwdKeeperConfig{}, nil
		}
		return nil, err
	}
	return cfg,nil
}
*/

// 创建初始化标志文件
func CreateInitFlag() {
	fileutil.WriteFile(InitFlagFile, timeutil.GetCurrentFmtTime())
}

