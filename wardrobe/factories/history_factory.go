package factories

import (
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
)

func HistoryFactory() models.History {
	num := gofakeit.Number(2, 3)

	return models.History{
		HistoryType:    gofakeit.LoremIpsumSentence(num - 1),
		HistoryContext: gofakeit.LoremIpsumSentence(num),
	}
}
