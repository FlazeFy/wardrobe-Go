package repositories

import (
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wash Interface
type WashRepository interface {
	HardDeleteWashByClothesID(clothesID, createdBy uuid.UUID) error
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
func (r *washRepository) HardDeleteWashByClothesID(clothesID, createdBy uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("clothes_id = ? AND created_by = ?", clothesID, createdBy).Delete(&models.Wash{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

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
