package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserRequest struct {
		ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		RequestToken   string    `json:"request_token" gorm:"type:varchar(6);not null"`
		RequestContext string    `json:"request_context" gorm:"type:varchar(14);not null"`
		CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		ValidatedAt    time.Time `json:"validated_at" gorm:"type:timestamp;null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
