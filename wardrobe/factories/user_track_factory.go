package factories

import (
	"fmt"
	"math/rand"
	"time"
	"wardrobe/config"
	"wardrobe/models"
)

func UserTrackFactory() models.UserTrack {
	rand.Seed(time.Now().UnixNano())

	lat := -90 + rand.Float64()*180
	long := -180 + rand.Float64()*360
	trackSource := config.TrackSources[rand.Intn(len(config.TrackSources))]

	return models.UserTrack{
		TrackLat:    fmt.Sprintf("%f", lat),
		TrackLong:   fmt.Sprintf("%f", long),
		TrackSource: trackSource,
	}
}
