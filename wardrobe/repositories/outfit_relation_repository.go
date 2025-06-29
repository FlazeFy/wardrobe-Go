package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Outfit Relation Interface
type OutfitRelationRepository interface {
	CreateOutfitRelation(outfitRel *models.OutfitRelation, userID uuid.UUID) error
	HardDeleteOutfitRelationByClothesID(clothesID, createdBy uuid.UUID) error

	// For Seeder
	DeleteAll() error
}

// Outfit Relation Struct
type outfitRelationRepository struct {
	db *gorm.DB
}

// Outfit Relation Constructor
func NewOutfitRelationRepository(db *gorm.DB) OutfitRelationRepository {
	return &outfitRelationRepository{db: db}
}

func (r *outfitRelationRepository) CreateOutfitRelation(outfitRel *models.OutfitRelation, userID uuid.UUID) error {
	// Default Value
	outfitRel.ID = uuid.New()
	outfitRel.CreatedAt = time.Now()
	outfitRel.CreatedBy = userID

	// Query
	return r.db.Create(outfitRel).Error
}

func (r *outfitRelationRepository) HardDeleteOutfitRelationByClothesID(clothesID, createdBy uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("clothes_id = ? AND created_by = ?", clothesID, createdBy).Delete(&models.OutfitRelation{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// For Seeder
func (r *outfitRelationRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.OutfitRelation{}).Error
}
