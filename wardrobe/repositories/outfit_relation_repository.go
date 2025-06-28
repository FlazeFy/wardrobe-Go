package repositories

import (
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Outfit Relation Interface
type OutfitRelationRepository interface {
	HardDeleteOutfitRelationByClothesID(clothesID, createdBy uuid.UUID) error
}

// Outfit Relation Struct
type outfitRelationRepository struct {
	db *gorm.DB
}

// Outfit Relation Constructor
func NewOutfitRelationRepository(db *gorm.DB) OutfitRelationRepository {
	return &outfitRelationRepository{db: db}
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
