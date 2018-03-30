package timeutil

import (
	"time"
	"fmt"
)

/**
 * 获取当前格式化的时间.
 */
func GetCurrentFmtTime() string {
	t := time.Now()
	return fmt.Sprintf("%4d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// 获取今天的数字日期字符串，格式如：20180104
func GetTodayNumericDateString() string {
	t := time.Now()
	return fmt.Sprintf("%4d%02d%02d", t.Year(), t.Month(), t.Day())
}


// 将YYYY-MM-DD HH:II:SS的时间字符串转换为unix时间戳
func GetTimeFromString(fmtTime string) time.Time {
	tm, _ := time.Parse("2006-01-02 15:04:05", fmtTime)
	return tm
}

// 获取格式化时间
func GetFormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}