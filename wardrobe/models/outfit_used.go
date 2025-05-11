package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	OutfitUsed struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Outfit
		OutfitId uuid.UUID `json:"outfit_id" gorm:"not null"`
		Outfit   Outfit    `json:"-" gorm:"foreignKey:OutfitId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
