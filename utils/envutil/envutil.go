package envutil

import (
	"os/user"
	"os"
	"bytes"
	"os/exec"
	"strings"
	"errors"
	"runtime"
	"bufio"
	"github.com/bgentry/speakeasy"
)

// 获取home目录路径，需要区分windows以及linux
func GetHomeDir() (string,error){
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support
	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

// 获取linux系统的home目录
func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

// 获取windows系统的home目录
func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

// 读取一行
func ReadLine() string {
	bio := bufio.NewReader(os.Stdin)
	line, _, _ := bio.ReadLine()
	return strings.Trim(string(line),"\r\n ")
}

// 读取一个字符
func ReadChar() string {
	bio := bufio.NewReader(os.Stdin)
	b,_ := bio.ReadByte()
	return string(b)
}

// 读取密码
func ReadPassword(tip string) (string, error) {
	password, err := speakeasy.Ask(tip)
	if err != nil {
		return "",err
	}
	return password, nil
}
