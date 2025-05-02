package models

import (
	"gorm.io/gorm"
)

type (
	Admin struct {
		gorm.Model
		Username        string `json:"username" gorm:"type:varchar(36);not null"`
		Password        string `json:"password" gorm:"type:varchar(500);not null"`
		Email           string `json:"email" gorm:"type:varchar(500);not null"`
		TelegramUserId  int    `json:"telegram_user_id" gorm:"null"`
		TelegramIsValid bool   `json:"telegram_is_valid"`
	}
)
