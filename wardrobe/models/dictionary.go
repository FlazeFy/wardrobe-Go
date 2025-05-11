package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Dictionary struct {
		ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		DictionaryType string    `json:"dictionary_type" gorm:"type:varchar(36);not null"`
		DictionaryName string    `json:"dictionary_name" gorm:"type:varchar(75);unique;not null"`
		CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
)
