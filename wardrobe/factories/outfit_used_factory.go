package factories

import (
	"wardrobe/models"

	"github.com/google/uuid"
)

func OutfitUsedFactory(outfitID uuid.UUID) models.OutfitUsed {
	return models.OutfitUsed{
		OutfitId: outfitID,
	}
}
