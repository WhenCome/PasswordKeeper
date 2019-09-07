package backuputil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/whencome/PasswordKeeper/config"
	"github.com/whencome/PasswordKeeper/utils/fileutil"
)

// 同步备份配置
func Sync() error {
	// 读取配置
	pwdCfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	backupDir := pwdCfg.UserCfg.BackupDir
	exists, err := fileutil.IsPathExists(backupDir)
	if err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(backupDir, 0777)
		if err != nil {
			return err
		}
	}
	syncDir(config.AppDataDir, backupDir)
	return nil
}

// 同步文件
func syncFile(srcFile, targetFile string) {
	fileutil.CopyFile(targetFile, srcFile)
}

// 同步目录
func syncDir(dir string, targetDir string) {
	exists, err := fileutil.IsPathExists(targetDir)
	if err != nil {
		return
	}
	if !exists {
		err = os.MkdirAll(targetDir, 0777)
		if err != nil {
			return
		}
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 不同步目录本身
		if path == dir {
			return nil
		}
		if info.IsDir() {
			syncDir(fmt.Sprintf("%s/%s", dir, info.Name()), fmt.Sprintf("%s/%s", targetDir, info.Name()))
		} else {
			syncFile(fmt.Sprintf("%s/%s", dir, info.Name()), fmt.Sprintf("%s/%s", targetDir, info.Name()))
		}
		return nil
	})
}
