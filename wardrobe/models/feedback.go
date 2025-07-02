package models

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	FeedbackRate int       `json:"feedback_rate" gorm:"type:integer;not null" binding:"required,min=1,max=5"`
	FeedbackBody string    `json:"feedback_body" gorm:"type:varchar(255);not null" binding:"required,max=255"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	// FK - User
	CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
