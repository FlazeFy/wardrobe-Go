package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null" binding:"required,min=6,max=36"`
		Password        string    `json:"password" gorm:"type:varchar(500);not null" binding:"required,min=6,max=36"`
		Email           string    `json:"email" gorm:"type:varchar(500);not null" binding:"required,email,max=500"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null" binding:"omitempty,max=36"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	UserContact struct {
		Username        string  `json:"username"`
		Email           string  `json:"email"`
		TelegramUserId  *string `json:"telegram_user_id"`
		TelegramIsValid bool    `json:"telegram_is_valid"`
	}
	UserReadyFetchWeather struct {
		UserID          uuid.UUID `json:"user_id"`
		TrackLat        string    `json:"track_lat"`
		TrackLong       string    `json:"track_long"`
		CreatedAt       time.Time `json:"created_at"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
	}
)

// For Generic Interface
func (a *User) GetID() uuid.UUID {
	return a.ID
}
func (a *User) GetPassword() string {
	return a.Password
}
