package factories

import (
	"wardrobe/models"

	"github.com/google/uuid"
)

func OutfitRelationFactory(outfitID, clothesID uuid.UUID) models.OutfitRelation {
	return models.OutfitRelation{
		OutfitId:  outfitID,
		ClothesId: clothesID,
	}
}
