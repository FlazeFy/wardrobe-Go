package factories

import (
	"wardrobe/config"
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func ScheduleFactory(clothesID uuid.UUID) models.Schedule {
	scheduleNote := gofakeit.LoremIpsumSentence(gofakeit.Number(3, 10))

	return models.Schedule{
		Day:          gofakeit.RandomString(config.Days),
		ScheduleNote: &scheduleNote,
		IsRemind:     gofakeit.Bool(),
		ClothesId:    clothesID,
	}
}
