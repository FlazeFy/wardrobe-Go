package factories

import (
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
)

func AdminFactory() models.Admin {
	password := "nopass123"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return models.Admin{
		Username:        gofakeit.Username(),
		Password:        string(hashedPass),
		TelegramUserId:  nil,
		TelegramIsValid: false,
		Email:           gofakeit.Email(),
	}
}
