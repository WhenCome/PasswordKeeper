package backuputil

import (
	"../../config"
	"path/filepath"
	"os"
	"fmt"
)

// 同步备份配置
func Sync() {
	filepath.Walk(config.AppDataDir, func(path string, info os.FileInfo, err error) error {

		return nil
	})
}

func syncFile(srcFile, fileName string){

}

func syncDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			syncDir(fmt.Sprintf("%s/%s", dir, info.Name()))
		} else {
			syncFile(fmt.Sprintf("%s/%s", dir, info.Name()), info.Name())
		}
		return nil
	})
}