package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleContext struct {
	DB *gorm.DB
}

func NewScheduleContext(db *gorm.DB) *ScheduleContext {
	return &ScheduleContext{DB: db}
}

type (
	Schedule struct {
		ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Day          string    `json:"day" gorm:"type:varchar(3);not null"`
		ScheduleNote *string   `json:"schedule_note" gorm:"type:varchar(255);null"`
		IsRemind     bool      `json:"is_remind" gorm:"type:boolean;not null"`
		CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Clothes
		ClothesId uuid.UUID `json:"clothes_id" gorm:"not null"`
		Clothes   Clothes   `json:"-" gorm:"foreignKey:ClothesId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}

	// Queries Only
	ScheduleByDay struct {
		ClothesName     string    `json:"clothes_name"`
		Day             string    `json:"day"`
		ScheduleNote    *string   `json:"schedule_note"`
		ClothesImage    *string   `json:"clothes_image"`
		ClothesType     string    `json:"clothes_type"`
		ClothesCategory string    `json:"clothes_category"`
		ClothesId       uuid.UUID `json:"clothes_id"`
	}
)

func (c *ScheduleContext) GetScheduleByDay(day string, userId uuid.UUID) ([]ScheduleByDay, error) {
	// Model
	var data []ScheduleByDay

	// Query
	result := c.DB.
		Table("schedules").
		Select("clothes_name,day,schedule_note,clothes_image,clothes_type,clothes_category,clothes.id AS clothes_id").
		Joins("JOIN clothes ON clothes.id = schedules.clothes_id").
		Where("day = ? AND schedules.created_by = ?", day, userId).
		Scan(&data)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

// Command Scheduler
func (c *ScheduleContext) SchedulerDeleteSchedulehById(id uuid.UUID) (int64, error) {
	// Model
	var schedule Schedule

	// Query
	result := c.DB.Unscoped().Where("clothes_id", id).Delete(&schedule)

	// Response
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
