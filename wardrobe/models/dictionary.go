package models

import (
	"gorm.io/gorm"
)

type (
	Dictionary struct {
		gorm.Model
		DictionaryType string `json:"dictionary_type" gorm:"type:varchar(36);not null"`
		DictionaryName string `json:"dictionary_name" gorm:"type:varchar(75);unique;not null"`
	}
)
