package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	GoogleToken struct {
		ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		AccessToken  string    `json:"access_token" gorm:"type:text;not null"`
		RefreshToken string    `json:"refresh_token" gorm:"type:text;not null"`
		Expiry       time.Time `json:"expiry" gorm:"type:timestamp;not null"`
		CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
