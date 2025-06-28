package factories

import (
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/brianvoe/gofakeit/v6"
)

func ErrorFactory() models.Error {
	line := uint(utils.GetRandInt(500, 10))
	randIntStackTrace := utils.GetRandInt(25, 12)
	randIntMessage := utils.GetRandInt(15, 7)
	randIntFile := utils.GetRandInt(4, 2)

	return models.Error{
		Message:    gofakeit.LoremIpsumSentence(randIntMessage),
		StackTrace: gofakeit.LoremIpsumSentence(randIntStackTrace),
		File:       gofakeit.LoremIpsumSentence(randIntFile),
		Line:       line,
	}
}
