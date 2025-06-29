package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Outfit Interface
type OutfitRepository interface {
	CreateOutfit(outfit *models.Outfit, userID uuid.UUID) error

	// For Seeder
	DeleteAll() error
	FindOneRandom(userID uuid.UUID) (*models.Outfit, error)
}

// Outfit Struct
type outfitRepository struct {
	db *gorm.DB
}

// Outfit Constructor
func NewOutfitRepository(db *gorm.DB) OutfitRepository {
	return &outfitRepository{db: db}
}

func (r *outfitRepository) CreateOutfit(outfit *models.Outfit, userID uuid.UUID) error {
	// Default Value
	outfit.ID = uuid.New()
	outfit.CreatedAt = time.Now()
	outfit.CreatedBy = userID

	// Query
	return r.db.Create(outfit).Error
}

// For Seeder
func (r *outfitRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.Outfit{}).Error
}

func (r *outfitRepository) FindOneRandom(userID uuid.UUID) (*models.Outfit, error) {
	var outfit models.Outfit

	err := r.db.Where("created_by", userID).Order("RANDOM()").Limit(1).First(&outfit).Error
	if err != nil {
		return nil, err
	}

	return &outfit, nil
}
