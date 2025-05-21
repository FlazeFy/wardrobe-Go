package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

type WashContext struct {
	DB *gorm.DB
}

func NewWashContext(db *gorm.DB) *WashContext {
	return &WashContext{DB: db}
}

// Command Scheduler
func (c *WashContext) SchedulerDeleteWashById(id uuid.UUID) (int64, error) {
	// Model
	var wash Wash

	// Query
	result := c.DB.Unscoped().Where("clothes_id", id).Delete(&wash)

	// Response
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
