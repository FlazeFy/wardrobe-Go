package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Outfit Used Interface
type OutfitUsedRepository interface {
	CreateOutfitUsed(outfitUsed *models.OutfitUsed, userID uuid.UUID) error

	// For Seeder
	DeleteAll() error
}

// Outfit Used Struct
type outfitUsedRepository struct {
	db *gorm.DB
}

// Outfit Used Constructor
func NewOutfitUsedRepository(db *gorm.DB) OutfitUsedRepository {
	return &outfitUsedRepository{db: db}
}

func (r *outfitUsedRepository) CreateOutfitUsed(outfitUsed *models.OutfitUsed, userID uuid.UUID) error {
	// Default Value
	outfitUsed.ID = uuid.New()
	outfitUsed.CreatedAt = time.Now()
	outfitUsed.CreatedBy = userID

	// Query
	return r.db.Create(outfitUsed).Error
}

// For Seeder
func (r *outfitUsedRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.OutfitUsed{}).Error
}
