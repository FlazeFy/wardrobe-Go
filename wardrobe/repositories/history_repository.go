package repositories

import (
	"time"
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// History Interface
type HistoryRepository interface {
	FindAllHistory(pagination utils.Pagination, userID *uuid.UUID) ([]models.GetHistory, int64, error)
	HardDeleteHistoryByID(ID, createdBy uuid.UUID) error
	CreateHistory(history *models.History, userID uuid.UUID) error
	// Task Scheduler
	DeleteHistoryForLastNDays(days int) (int64, error)
	// For Seeder
	DeleteAll() error
}

// History Struct
type historyRepository struct {
	db *gorm.DB
}

// History Constructor
func NewHistoryRepository(db *gorm.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) FindAllHistory(pagination utils.Pagination, userID *uuid.UUID) ([]models.GetHistory, int64, error) {
	// Model
	var histories []models.GetHistory
	var total int64

	// Pagination Count
	offset := (pagination.Page - 1) * pagination.Limit
	countQuery := r.db.Table("histories").
		Joins("JOIN users ON users.id = histories.created_by")
	if userID != nil {
		countQuery = countQuery.Where("created_by = ?", userID)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Query
	query := r.db.Table("histories").
		Select("histories.id, history_type, history_context, histories.created_at, username").
		Joins("JOIN users ON users.id = histories.created_by").
		Limit(pagination.Limit).
		Offset(offset)
	if userID != nil {
		query = query.Where("created_by = ?", userID)
	}
	if err := query.Find(&histories).Error; err != nil {
		return nil, 0, err
	}
	if len(histories) == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	return histories, total, nil
}

func (r *historyRepository) CreateHistory(history *models.History, userID uuid.UUID) error {
	// Default
	history.ID = uuid.New()
	history.CreatedAt = time.Now()
	history.CreatedBy = userID

	// Query
	return r.db.Create(history).Error
}

func (r *historyRepository) HardDeleteHistoryByID(ID, createdBy uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("id = ?", ID).Where("created_by = ?", createdBy).Delete(&models.History{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Task Scheduler
func (r *historyRepository) DeleteHistoryForLastNDays(days int) (int64, error) {
	// Cutoff Days
	cutoff := time.Now().AddDate(0, 0, -days)

	// Query
	result := r.db.Unscoped().Where("created_at < ?", cutoff).Delete(&models.History{})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// For Seeder
func (r *historyRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.History{}).Error
}
