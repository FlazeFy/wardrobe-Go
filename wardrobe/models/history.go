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
		ID             uuid.UUID `json:"id" example:"a01bd5b9-6cab-48a6-bec3-cb19fe07372e"`
		HistoryType    string    `json:"history_type" example:"Quam rerum."`
		HistoryContext string    `json:"history_context" example:"Voluptatibus nihil accusantium."`
		CreatedAt      time.Time `json:"created_at" example:"a01bd5b9-6cab-48a6-bec3-cb19fe07372e"`
		Username       string    `json:"username" example:"flazefy"`
	}
)
