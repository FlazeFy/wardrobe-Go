package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Error struct {
		ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Message    string    `json:"message" gorm:"type:text;not null"`
		StackTrace string    `json:"stack_trace" gorm:"type:text;not null"`
		File       string    `json:"file" gorm:"type:varchar(255);not null"`
		Line       uint      `json:"line" gorm:"not null"`
		IsFixed    bool      `json:"is_fixed" gorm:"not null"`
		CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	ErrorAudit struct {
		Message   string `json:"message"`
		CreatedAt string `json:"created_at"`
		Total     int    `json:"total"`
	}
)
