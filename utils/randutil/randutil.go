package randutil

import (
	"time"
	"fmt"
	"strings"
	"math/rand"
	"bytes"
)

var Alphabets []string = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var Digits []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

// 获取指定长度的随机数字序列
func GetRandDigitString(length int) string {
	rand.Seed(time.Now().UnixNano())
	s := fmt.Sprintf("%d", rand.Uint32())
	if len(s) >= length {
		return s[len(s)-length:]
	}
	return strings.Repeat("0", length - len(s)) + s
}

// 生成给定范围的随机数
func GenerateRangeRandNumber(min, max int) int {
	randNum := rand.Intn(max - min) + min
	return randNum
}

// 获取随机的字符串，只包含字母和数字
func GetRandAlphaDigitString(length int) string {
	var buffer bytes.Buffer
	var randChars []string
	randChars = append(randChars, Alphabets...)
	randChars = append(randChars, Digits...)
	rand.Seed(time.Now().UnixNano())
	var idx int = 0
	randPos := 0
	for idx<length {
		randPos = GenerateRangeRandNumber(0, len(randChars)-1)
		buffer.WriteString(randChars[randPos])
		idx++
	}
	return buffer.String()
}

