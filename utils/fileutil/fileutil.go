package fileutil

import (
	"os"
	"strings"
	"fmt"
)

// 检查文件是否存在
func IsFileExists(file string) bool {
	finfo, err := os.Stat(file)
	if err == nil && !finfo.IsDir() {
		return true
	}
	return false
}

// 检查给定路径是否存在
func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 构造路径
func BuildPath(path, fileName string) string {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		return fmt.Sprintf("%s%s", path, fileName)
	}
	return fmt.Sprintf("%s/%s", path, fileName)
}
