package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Admin struct {
		ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		Password        string    `json:"password" gorm:"type:varchar(500);not null"`
		Email           string    `json:"email" gorm:"type:varchar(500);not null"`
		TelegramUserId  int       `json:"telegram_user_id" gorm:"null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
)
