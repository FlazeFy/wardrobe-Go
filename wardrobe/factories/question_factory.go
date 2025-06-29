package factories

import (
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
)

func QuestionFactory() models.Question {
	len := gofakeit.Number(4, 15)

	return models.Question{
		Question: gofakeit.LoremIpsumSentence(len),
	}
}
