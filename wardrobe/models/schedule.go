package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Schedule struct {
		ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Day          string    `json:"day" gorm:"type:varchar(3);not null" binding:"required,min=3,max=3"`
		ScheduleNote *string   `json:"schedule_note" gorm:"type:varchar(255);null" binding:"omitempty,max=255"`
		IsRemind     bool      `json:"is_remind" gorm:"type:boolean;not null" binding:"required"`
		CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Clothes
		ClothesId uuid.UUID `json:"clothes_id" gorm:"not null" binding:"required,max=36,min=36"`
		Clothes   Clothes   `json:"-" gorm:"foreignKey:ClothesId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}

	// Queries Only
	ScheduleByDay struct {
		ClothesName     string    `json:"clothes_name"`
		Day             string    `json:"day"`
		ScheduleNote    *string   `json:"schedule_note"`
		ClothesImage    *string   `json:"clothes_image"`
		ClothesType     string    `json:"clothes_type"`
		ClothesCategory string    `json:"clothes_category"`
		ClothesId       uuid.UUID `json:"clothes_id"`
	}
	ScheduleReadyToCalendarTask struct {
		ID           uuid.UUID `json:"id"`
		ClothesName  string    `json:"clothes_name"`
		Day          string    `json:"day"`
		ScheduleNote *string   `json:"schedule_note"`
		Username     string    `json:"username"`
		AccessToken  string    `json:"access_token"`
	}
)
