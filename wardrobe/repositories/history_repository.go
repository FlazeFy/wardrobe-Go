package repositories

import (
	"time"
	"wardrobe/models"

	"gorm.io/gorm"
)

// History Interface
type HistoryRepository interface {
	// Task Scheduler
	DeleteHistoryForLastNDays(days int) (int64, error)
}

// History Struct
type historyRepository struct {
	db *gorm.DB
}

// History Constructor
func NewHistoryRepository(db *gorm.DB) HistoryRepository {
	return &historyRepository{db: db}
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
