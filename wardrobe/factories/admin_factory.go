package factories

import (
	"wardrobe/models"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
)

func AdminFactory(username, email, telegramUserId, password *string) models.Admin {
	var finalUsername string
	if username != nil && *username != "" {
		finalUsername = *username
	} else {
		finalUsername = gofakeit.Username()
	}

	var finalEmail string
	if email != nil && *email != "" {
		finalEmail = *email
	} else {
		finalEmail = gofakeit.Email()
	}

	var pwd string
	if password != nil && *password != "" {
		pwd = *password
	} else {
		pwd = "nopass123"
	}
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	var finalTelegramUserId *string
	if telegramUserId != nil && *telegramUserId != "" {
		finalTelegramUserId = telegramUserId
	} else {
		finalTelegramUserId = nil
	}

	return models.Admin{
		Username:        finalUsername,
		Password:        string(hashedPass),
		TelegramUserId:  finalTelegramUserId,
		TelegramIsValid: false,
		Email:           finalEmail,
	}
}
