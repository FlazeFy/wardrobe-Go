package factories

import (
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
)

func FeedbackFactory() models.Feedback {
	rate := gofakeit.Number(1, 5)

	return models.Feedback{
		FeedbackRate: rate,
		FeedbackBody: gofakeit.LoremIpsumSentence(rate),
	}
}
