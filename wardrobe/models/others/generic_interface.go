package others

import (
	"time"

	"github.com/google/uuid"
)

type (
	GoogleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	Metadata struct {
		Total      int `json:"total"`
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		TotalPages int `json:"total_pages"`
	}
	// All Role
	Account interface {
		GetID() uuid.UUID
		GetPassword() string
	}
	MyProfile struct {
		Username        string    `json:"username" gorm:"type:varchar(36);not null" example:"testeruser"`
		Email           string    `json:"email" gorm:"type:varchar(500);not null" example:"testeruser@gmail.com"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null" example:"1317453312"`
		TelegramIsValid bool      `json:"telegram_is_valid" example:"false"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null" example:"2025-07-16T20:25:19.945914+07:00"`
	}
)
