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
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		Email           string    `json:"email" gorm:"type:varchar(500);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
)
