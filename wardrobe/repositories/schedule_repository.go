package repositories

import (
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schedule Interface
type ScheduleRepository interface {
	GetScheduleByDay(day string, userId uuid.UUID) ([]models.ScheduleByDay, error)
	SchedulerDeleteSchedulehById(id uuid.UUID) (int64, error)
}

// Schedule Struct
type scheduleRepository struct {
	db *gorm.DB
}

// Schedule Constructor
func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) GetScheduleByDay(day string, userId uuid.UUID) ([]models.ScheduleByDay, error) {
	// Model
	var data []models.ScheduleByDay

	// Query
	result := r.db.Table("schedules").
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
func (r *scheduleRepository) SchedulerDeleteSchedulehById(id uuid.UUID) (int64, error) {
	// Model
	var schedule models.Schedule

	// Query
	result := r.db.Unscoped().Where("clothes_id", id).Delete(&schedule)

	// Response
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
