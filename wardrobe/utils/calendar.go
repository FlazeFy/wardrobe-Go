package utils

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func AddWeeklyGoogleCalendarEvent(accessToken string, summary string, description string, startTime time.Time) (*calendar.Event, error) {
	ctx := context.Background()

	srv, err := calendar.NewService(ctx, option.WithTokenSource(
		oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: accessToken,
		}),
	))
	if err != nil {
		return nil, err
	}

	event := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
		End: &calendar.EventDateTime{
			DateTime: startTime.Add(1 * time.Hour).Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
		Recurrence: []string{
			"RRULE:FREQ=WEEKLY",
		},
	}

	createdEvent, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}
