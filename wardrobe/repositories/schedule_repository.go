package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schedule Interface
type ScheduleRepository interface {
	CheckScheduleByDayAndClothesID(day string, userId, clothesID uuid.UUID) (bool, error)
	CreateSchedule(schedule *models.Schedule, userID uuid.UUID) error
	FindScheduleByDay(day string, userId uuid.UUID) ([]models.ScheduleByDay, error)
	DeleteScheduleByClothesId(id uuid.UUID) (int64, error)
	HardDeleteScheduleById(id, createdBy uuid.UUID) error
	HardDeleteScheduleByClothesID(clothesID, createdBy uuid.UUID) error
}

// Schedule Struct
type scheduleRepository struct {
	db *gorm.DB
}

// Schedule Constructor
func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) FindScheduleByDay(day string, userId uuid.UUID) ([]models.ScheduleByDay, error) {
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
	if len(data) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, nil
}

func (r *scheduleRepository) CheckScheduleByDayAndClothesID(day string, userID, clothesID uuid.UUID) (bool, error) {
	// Model
	var data []models.ScheduleByDay

	// Query
	result := r.db.Table("schedules").
		Where("day = ? AND created_by = ? AND clothes_id = ?", day, userID, clothesID).
		First(&data)

	// Response
	if result.Error != nil {
		return true, result.Error
	}
	if len(data) == 0 {
		return false, nil
	}

	return true, nil
}

func (r *scheduleRepository) CreateSchedule(schedule *models.Schedule, userID uuid.UUID) error {
	// Default
	schedule.ID = uuid.New()
	schedule.CreatedAt = time.Now()
	schedule.CreatedBy = userID

	// Query
	return r.db.Create(schedule).Error
}

func (r *scheduleRepository) HardDeleteScheduleByClothesID(clothesID, createdBy uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("clothes_id = ? AND created_by = ?", clothesID, createdBy).Delete(&models.Schedule{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// fix this to used user id also.
func (r *scheduleRepository) DeleteScheduleByClothesId(id uuid.UUID) (int64, error) {
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

func (r *scheduleRepository) HardDeleteScheduleById(id, createdBy uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("id", id).Where("created_by", createdBy).Delete(&models.Schedule{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
