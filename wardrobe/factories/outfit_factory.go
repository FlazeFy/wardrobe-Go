package factories

import (
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
)

func OutfitFactory() models.Outfit {
	outfitNote := gofakeit.LoremIpsumSentence(gofakeit.Number(1, 3))

	return models.Outfit{
		OutfitName: gofakeit.LoremIpsumSentence(gofakeit.Number(1, 2)),
		OutfitNote: &outfitNote,
		IsAuto:     gofakeit.Bool(),
		IsFavorite: gofakeit.Bool(),
	}
}
