package utils

import "time"

func GetFormattedDate(date int64) string {
	return time.Unix(date, 0).Format("2006-01-02:15:04:05")
}
