package utils

import (
	"time"
)

func GetNextDay(current string, offset int) string {
	now := time.Now()
	targetDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for i := 0; i < 7; i++ {
		targetDay = targetDay.AddDate(0, 0, 1)
		if targetDay.Weekday().String()[:3] == current {
			targetDay = targetDay.AddDate(0, 0, offset)
			break
		}
	}

	return targetDay.Weekday().String()[:3]
}
