package unittest

import (
	"testing"
	"time"
	"wardrobe/utils"
)

func GeneratorTestGetNextDay(t *testing.T) {
	type args struct {
		current string
		offset  int
	}

	// Sample Test
	tests := []struct {
		name        string
		args        args
		wantWeekday string
	}{
		{
			name: "Next Tuesday + 0 offset",
			args: args{
				current: "Tue",
				offset:  0,
			},
			wantWeekday: "Tue",
		},
		{
			name: "Next Wednesday + 2 offset",
			args: args{
				current: "Wed",
				offset:  2,
			},
			wantWeekday: time.Now().AddDate(0, 0, daysUntil("Wed")+2).Weekday().String()[:3],
		},
		{
			name: "Next Sunday + 1 offset",
			args: args{
				current: "Sun",
				offset:  1,
			},
			wantWeekday: time.Now().AddDate(0, 0, daysUntil("Sun")+1).Weekday().String()[:3],
		},
	}

	// Exec
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.GetNextDay(tt.args.current, tt.args.offset)
			if got != tt.wantWeekday {
				t.Errorf("GetNextDay(%v, %d) = %v, want %v", tt.args.current, tt.args.offset, got, tt.wantWeekday)
			}
		})
	}
}

func daysUntil(target string) int {
	now := time.Now()
	for i := 1; i <= 7; i++ {
		next := now.AddDate(0, 0, i)
		if next.Weekday().String()[:3] == target {
			return i
		}
	}
	return 0
}
