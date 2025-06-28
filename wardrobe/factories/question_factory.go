package factories

import (
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/brianvoe/gofakeit/v6"
)

func QuestionFactory() models.Question {
	len := utils.GetRandInt(15, 4)

	return models.Question{
		Question: gofakeit.LoremIpsumSentence(len),
	}
}
