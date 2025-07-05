package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Question struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Question  string    `json:"question" gorm:"type:varchar(500);not null" binding:"required,min=10,max=500"`
		Answer    *string   `json:"answer" gorm:"type:varchar(500);null" binding:"omitempty,max=500"`
		IsShow    bool      `json:"is_show" gorm:"type:boolean;not null"`
		CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	UnansweredQuestion struct {
		Question  string    `json:"question" gorm:"type:varchar(500);not null"`
		CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
)
