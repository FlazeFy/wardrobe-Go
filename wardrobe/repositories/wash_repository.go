package repositories

import (
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wash Interface
type WashRepository interface {
	DeleteWashByClothesId(id uuid.UUID) (int64, error)
}

// Wash Struct
type washRepository struct {
	db *gorm.DB
}

// Wash Constructor
func NewWashRepository(db *gorm.DB) WashRepository {
	return &washRepository{db: db}
}

// Command Scheduler
func (r *washRepository) DeleteWashByClothesId(id uuid.UUID) (int64, error) {
	// Model
	var wash models.Wash

	// Query
	result := r.db.Unscoped().Where("clothes_id", id).Delete(&wash)

	// Response
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
