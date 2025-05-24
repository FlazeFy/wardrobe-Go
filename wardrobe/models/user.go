package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserContext struct {
	DB *gorm.DB
}

func NewUserContext(db *gorm.DB) *UserContext {
	return &UserContext{DB: db}
}

type (
	User struct {
		ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		Password        string    `json:"password" gorm:"type:varchar(500);not null"`
		Email           string    `json:"email" gorm:"type:varchar(500);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
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

func (c *UserContext) SchedulerGetUserReadyFetchWeather() ([]UserReadyFetchWeather, error) {
	// Model
	var users []UserReadyFetchWeather

	// Query
	result := c.DB.Table("users").
		Select(`
			user_tracks.track_lat,user_tracks.track_long,user_tracks.created_at,users.id as user_id,
			users.username,users.telegram_user_id,users.telegram_is_valid
		`).
		Joins(`
			JOIN (
				SELECT DISTINCT ON (created_by) *
				FROM user_tracks
				ORDER BY created_by, created_at DESC
			) AS user_tracks ON user_tracks.created_by = users.id
		`).
		Where("user_tracks.track_lat IS NOT NULL AND user_tracks.track_long IS NOT NULL").
		Order("users.username ASC").
		Find(&users)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}
	if len(users) == 0 {
		return nil, errors.New("no user track found")
	}

	return users, nil
}
