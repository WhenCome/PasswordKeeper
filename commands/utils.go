package commands

import "strings"

// 包装字符串，便于输出
func wrapString(input string, length int) string {
	if len(input) > length {
		return input[:length]
	}
	return input+strings.Repeat(" ", length - len(input))
}
