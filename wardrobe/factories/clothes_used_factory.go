package factories

import (
	"wardrobe/config"
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func ClothesUsedFactory(clothesID uuid.UUID) models.ClothesUsed {
	clothesNote := gofakeit.Sentence(gofakeit.Number(3, 10))

	return models.ClothesUsed{
		ClothesNote: &clothesNote,
		UsedContext: gofakeit.RandomString(config.UsedContexts),
		ClothesId:   clothesID,
	}
}
