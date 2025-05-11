package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	ClothesUsed struct {
		ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesNote *string   `json:"clothes_desc" gorm:"type:varchar(500);null"`
		CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Clothes
		ClothesId uuid.UUID `json:"clothes_id" gorm:"not null"`
		Clothes   Clothes   `json:"-" gorm:"foreignKey:ClothesId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Used Context
		UsedContext uuid.UUID  `json:"used_context" gorm:"not null"`
		Dictionary  Dictionary `json:"-" gorm:"foreignKey:UsedContext;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
