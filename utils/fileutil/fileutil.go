package fileutil

import (
	"os"
	"strings"
	"fmt"
	"io"
	"io/ioutil"
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

// 复制文件
func CopyFile(dstName, srcName string) (int64, error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// 读取文件内容
func GetContents(filePath string) ([]byte,error) {
	f,err := os.Open(filePath)
	if err != nil {
		return nil,err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// 写入文件
func WriteFile(filePath, content string) (int, error) {
	fmt.Printf("File: %s;  Content: %s \n", filePath, content)
	var f *os.File
	var err error
	if !IsFileExists(filePath) {
		f, err = os.Create(filePath)
	} else {
		f, err = os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0666)
	}
	if err != nil {
		return 0,err
	}
	defer f.Close()
	return f.WriteString(content)
}
