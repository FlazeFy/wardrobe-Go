package factories

import (
	"wardrobe/config"
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
)

func ClothesFactory() models.Clothes {
	desc := gofakeit.Sentence(gofakeit.Number(3, 10))
	merk := gofakeit.Company()
	price := gofakeit.Number(100000, 2000000)
	buyAt := gofakeit.Date()

	return models.Clothes{
		ClothesName:     gofakeit.ProductName(),
		ClothesDesc:     &desc,
		ClothesMerk:     &merk,
		ClothesColor:    gofakeit.RandomString(config.Colors),
		ClothesPrice:    &price,
		ClothesBuyAt:    &buyAt,
		ClothesQty:      gofakeit.Number(1, 20),
		IsFaded:         gofakeit.Bool(),
		HasWashed:       gofakeit.Bool(),
		HasIroned:       gofakeit.Bool(),
		IsFavorite:      gofakeit.Bool(),
		IsScheduled:     gofakeit.Bool(),
		ClothesMadeFrom: gofakeit.RandomString(config.ClothesMadeFroms),
		ClothesType:     gofakeit.RandomString(config.ClothesTypes),
		ClothesCategory: gofakeit.RandomString(config.ClothesCategories),
		ClothesSize:     gofakeit.RandomString(config.ClothesSizes),
		ClothesGender:   gofakeit.RandomString(config.ClothesGenders),
	}
}
