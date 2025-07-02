package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	ClothesUsed struct {
		ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesNote *string   `json:"clothes_note" gorm:"type:varchar(500);null" binding:"omitempty,max=500"`
		CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Clothes
		ClothesId uuid.UUID `json:"clothes_id" gorm:"not null" binding:"required,max=36,min=36"`
		Clothes   Clothes   `json:"-" gorm:"foreignKey:ClothesId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Used Context
		UsedContext string     `json:"used_context" gorm:"not null" binding:"required,min=36"`
		Dictionary  Dictionary `json:"-" gorm:"foreignKey:UsedContext;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	ClothesUsedHistory struct {
		ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesName string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		ClothesType string    `json:"clothes_type" gorm:"type:varchar(36);not null"`
		ClothesNote *string   `json:"clothes_note" gorm:"type:varchar(500);null"`
		CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - Used Context
		UsedContext string `json:"used_context" gorm:"not null"`
	}
	SchedulerUsedClothesReadyToWash struct {
		ClothesName     string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		ClothesType     string    `json:"clothes_type" gorm:"type:varchar(36);not null"`
		ClothesMadeFrom string    `json:"clothes_made_from" gorm:"type:varchar(36);not null"`
		IsFaded         bool      `json:"is_faded" gorm:"type:boolean;not null"`
		IsScheduled     bool      `json:"is_scheduled" gorm:"type:boolean;not null"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		// FK - Used Context
		UsedContext string `json:"used_context" gorm:"not null"`
	}
)
