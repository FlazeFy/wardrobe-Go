package services

import (
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// Wash Interface
type WashService interface {
	DeleteWashByClothesId(id uuid.UUID) (int64, error)
}

// Wash Struct
type washService struct {
	washRepo repositories.WashRepository
}

// Wash Constructor
func NewWashService(washRepo repositories.WashRepository) WashService {
	return &washService{
		washRepo: washRepo,
	}
}

func (s *washService) DeleteWashByClothesId(id uuid.UUID) (int64, error) {
	// Repo : Delete Wash By Clothes Id
	rows, err := s.washRepo.DeleteWashByClothesId(id)
	if err != nil {
		return 0, err
	}

	return rows, nil
}
