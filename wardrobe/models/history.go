package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	History struct {
		ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		HistoryType    string    `json:"history_type" gorm:"type:varchar(36);not null"`
		HistoryContext string    `json:"history_context" gorm:"type:varchar(255);not null"`
		CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	GetHistory struct {
		ID             uuid.UUID `json:"id"`
		HistoryType    string    `json:"history_type"`
		HistoryContext string    `json:"history_context"`
		CreatedAt      time.Time `json:"created_at"`
		Username       string    `json:"username"`
	}
)
