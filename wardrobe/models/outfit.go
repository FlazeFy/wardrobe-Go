package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Outfit struct {
		ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		OutfitName string    `json:"outfit_name" gorm:"type:varchar(36);not null" binding:"required,max=36"`
		OutfitNote *string   `json:"outfit_note" gorm:"type:varchar(255);null" binding:"omitempty,max=255"`
		IsAuto     bool      `json:"is_auto" gorm:"type:boolean;not null" binding:"required"`
		IsFavorite bool      `json:"is_favorite" gorm:"type:boolean;not null" binding:"required"`
		CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		UpdatedAt  time.Time `json:"updated_at" gorm:"type:timestamp;null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
