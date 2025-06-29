package utils

import (
	"math"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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

func GetRandWeatherTemp(min, max float64) float64 {
	raw := gofakeit.Float64Range(min, max)
	temp := math.Round(raw*100) / 100

	return temp
}
