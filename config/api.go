package config

import (
	"encoding/xml"
	"log"

	"github.com/whencome/PasswordKeeper/utils/fileutil"
	"github.com/whencome/PasswordKeeper/utils/timeutil"
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
	if AppConfig != nil {
		return AppConfig, nil
	} else {
		// 检查配置文件是否存在
		if !fileutil.IsFileExists(AppConfigFile) {
			return nil, ErrConfigNotExists
		}
		// 读取文件内容
		cfgContent, err := fileutil.GetContents(AppConfigFile)
		if err != nil {
			return nil, err
		}
		pwdCfg := &PwdKeeperConfig{}
		err = xml.Unmarshal(cfgContent, pwdCfg)
		if err != nil {
			return nil, err
		}
		AppConfig = pwdCfg
		return pwdCfg, nil
	}
}

// 创建初始化标志文件
func CreateInitFlag() {
	fileutil.WriteFile(InitFlagFile, timeutil.GetCurrentFmtTime())
}

// 获取私钥
func GetPrivateKey() []byte {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Load config failed : %s", err)
	}
	privateKeyFile := cfg.CertCfg.PrivateKeyFile
	if !fileutil.IsFileExists(privateKeyFile) {
		log.Fatalln("Can not find private cert file!")
	}
	privateKey, err := fileutil.GetContents(privateKeyFile)
	if err != nil {
		log.Fatalln("Read private key failed!")
	}
	return privateKey
}

// 获取公钥
func GetPublicKey() []byte {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Load config failed : %s", err)
	}
	publicKeyFile := cfg.CertCfg.PublicKeyFile
	if !fileutil.IsFileExists(publicKeyFile) {
		log.Fatalln("Can not find public cert file!")
	}
	publicKey, err := fileutil.GetContents(publicKeyFile)
	if err != nil {
		log.Fatalln("Read public key failed!")
	}
	return publicKey
}
