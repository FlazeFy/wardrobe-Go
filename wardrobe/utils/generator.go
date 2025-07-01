package utils

import (
	"fmt"
	"math"
	"strings"
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

func GetTomorrowDayName() string {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return tomorrow.Weekday().String()[:3]
}

func GetThisWeekdayWithHour(dayStr string, hour, min int) (time.Time, error) {
	dayStr = strings.ToLower(dayStr)
	dayMap := map[string]time.Weekday{
		"sun": time.Sunday,
		"mon": time.Monday,
		"tue": time.Tuesday,
		"wed": time.Wednesday,
		"thu": time.Thursday,
		"fri": time.Friday,
		"sat": time.Saturday,
	}

	targetWeekday, ok := dayMap[dayStr]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid weekday string: %s", dayStr)
	}

	today := time.Now()
	daysUntil := (int(targetWeekday) - int(today.Weekday()) + 7) % 7
	targetDate := today.AddDate(0, 0, daysUntil)

	return time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), hour, min, 0, 0, time.Local), nil
}
