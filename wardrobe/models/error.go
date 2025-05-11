package models

import (
	"time"

	"gorm.io/gorm"
)

type Error struct {
	gorm.Model
	Message    string    `gorm:"type:text;not null"`
	StackTrace string    `gorm:"type:text;not null"`
	File       string    `gorm:"type:varchar(255);not null"`
	Line       uint      `gorm:"not null"`
	IsFixed    bool      `gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;not null"`
}
