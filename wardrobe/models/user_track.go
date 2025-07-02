package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserTrack struct {
		ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		TrackLat    string    `json:"track_lat" gorm:"type:varchar(255);not null" binding:"required,max=255"`
		TrackLong   string    `json:"track_long" gorm:"type:varchar(255);not null" binding:"required,max=255"`
		TrackSource string    `json:"track_source" gorm:"type:varchar(16);not null" binding:"required,max=16"`
		CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
