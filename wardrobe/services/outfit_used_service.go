package services

import (
	"wardrobe/repositories"
)

// Outfit Used Interface
type OutfitUsedService interface {
}

// Outfit Used Struct
type outfitUsedService struct {
	outfitUsedRepo repositories.OutfitUsedRepository
}

// Outfit Used Constructor
func NewOutfitUsedService(outfitUsedRepo repositories.OutfitUsedRepository) OutfitUsedService {
	return &outfitUsedService{
		outfitUsedRepo: outfitUsedRepo,
	}
}
