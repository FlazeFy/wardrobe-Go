package factories

import (
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
)

func ErrorFactory() models.Error {
	line := uint(gofakeit.Number(10, 500))
	randIntStackTrace := gofakeit.Number(12, 25)
	randIntMessage := gofakeit.Number(7, 15)
	randIntFile := gofakeit.Number(2, 4)

	return models.Error{
		Message:    gofakeit.LoremIpsumSentence(randIntMessage),
		StackTrace: gofakeit.LoremIpsumSentence(randIntStackTrace),
		File:       gofakeit.LoremIpsumSentence(randIntFile),
		Line:       line,
	}
}
