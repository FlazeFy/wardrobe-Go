package factories

import (
	"encoding/json"
	"time"
	"wardrobe/config"
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func WashFactory(clothesID uuid.UUID) models.Wash {
	washNote := gofakeit.LoremIpsumSentence(gofakeit.Number(3, 6))

	var finishedAt *time.Time
	if gofakeit.Bool() {
		t := gofakeit.DateRange(time.Now().Add(time.Hour), time.Now().AddDate(0, 0, gofakeit.Number(1, 2)))
		finishedAt = &t
	}

	var washCheckpoint []byte
	if gofakeit.Bool() {
		checkpoints := []models.WashCheckpoint{
			{CheckpointName: "Soak", IsFinished: gofakeit.Bool()},
			{CheckpointName: "Rinse", IsFinished: gofakeit.Bool()},
			{CheckpointName: "Dry", IsFinished: gofakeit.Bool()},
			{CheckpointName: "Ironed", IsFinished: gofakeit.Bool()},
		}
		jsonData, _ := json.Marshal(checkpoints)
		washCheckpoint = jsonData
	} else {
		washCheckpoint = nil
	}

	return models.Wash{
		WashNote:       &washNote,
		WashCheckpoint: washCheckpoint,
		WashType:       gofakeit.RandomString(config.WashTypes),
		FinishedAt:     finishedAt,
		ClothesId:      clothesID,
	}
}
