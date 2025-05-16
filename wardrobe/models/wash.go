package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type (
	Wash struct {
		ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
		WashNote       *string        `json:"wash_note" gorm:"type:varchar(75);null"`
		WashCheckpoint datatypes.JSON `json:"wash_checkpoint" gorm:"type:json;null"`
		FinishedAt     *time.Time     `json:"finished_at" gorm:"null"`
		CreatedAt      time.Time      `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Clothes
		ClothesId uuid.UUID `json:"clothes_id" gorm:"not null"`
		Clothes   Clothes   `json:"-" gorm:"foreignKey:ClothesId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Wash Type
		WashType   string     `json:"wash_type" gorm:"not null"`
		Dictionary Dictionary `json:"-" gorm:"foreignKey:WashType;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
