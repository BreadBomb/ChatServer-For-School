package util

import "time"

func GetCurrentTimestamp() string {
	return time.Now().Format("20060102150405")
}