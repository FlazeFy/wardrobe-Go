package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Dictionary struct {
		ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		DictionaryType string    `json:"dictionary_type" gorm:"type:varchar(36);not null" binding:"required,max=36"`
		DictionaryName string    `json:"dictionary_name" gorm:"type:varchar(75);unique;not null" binding:"required,min=3,max=75"`
		CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
)
