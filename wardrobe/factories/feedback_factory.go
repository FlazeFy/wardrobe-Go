package factories

import (
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/brianvoe/gofakeit/v6"
)

func FeedbackFactory() models.Feedback {
	rate := utils.GetRandInt(5, 1)

	return models.Feedback{
		FeedbackRate: rate,
		FeedbackBody: gofakeit.LoremIpsumSentence(rate),
	}
}
