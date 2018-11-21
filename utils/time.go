package utils

import "time"

const KMicTimeFormat = "2006/01/02 15:04:05.000000"

// Get a formatted Microseconds time
func GetMicTimeFormat() string {
	return time.Now().Format(KMicTimeFormat)
}

func GetTomorrowStartTime() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	return t.Unix() + 1
}

func GetNowTime() int64 {
	return time.Now().Unix()
}