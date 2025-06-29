package factories

import (
	"wardrobe/config"
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
)

func ClothesFactory() models.Clothes {
	clothesDesc := gofakeit.Sentence(gofakeit.Number(3, 10))
	clothesMerk := gofakeit.Company()
	clothesPrice := gofakeit.Number(100000, 2000000)
	clothesBuyAt := gofakeit.Date()

	return models.Clothes{
		ClothesName:     gofakeit.ProductName(),
		ClothesDesc:     &clothesDesc,
		ClothesMerk:     &clothesMerk,
		ClothesColor:    gofakeit.RandomString(config.Colors),
		ClothesPrice:    &clothesPrice,
		ClothesBuyAt:    &clothesBuyAt,
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
