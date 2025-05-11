package models

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	FeedbackRate int       `gorm:"type:integer;not null"`
	FeedbackBody string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	// FK - User
	CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
	User      User      `json:"-" gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
