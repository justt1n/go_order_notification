package main

import (
	"time"
)

// getCurrentTimeInGMT7 returns the current time in GMT+7 (Asia/Bangkok timezone)
func getCurrentTimeInGMT7() string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	currentTime := time.Now().In(loc)
	return currentTime.Format("2006-01-02 15:04:05")
}
