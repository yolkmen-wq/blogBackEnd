package utils

import "time"

// 用于解析时间的辅助函数
func ParseTime(createTimeBytes []byte) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", string(createTimeBytes))
}
