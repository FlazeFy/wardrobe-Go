package services

import (
	"wardrobe/models"
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// Clothes Used Interface
type ClothesUsedService interface {
	DeleteClothesUsedByClothesId(id uuid.UUID) (int64, error)

	// Task Scheduler
	SchedulerGetUsedClothesReadyToWash(days int) ([]models.SchedulerUsedClothesReadyToWash, error)
}

// Clothes Used Struct
type clothesUsedService struct {
	clothesUsedRepo repositories.ClothesUsedRepository
}

// Clothes Used Constructor
func NewClothesUsedService(clothesUsedRepo repositories.ClothesUsedRepository) ClothesUsedService {
	return &clothesUsedService{
		clothesUsedRepo: clothesUsedRepo,
	}
}

func (s *clothesUsedService) DeleteClothesUsedByClothesId(id uuid.UUID) (int64, error) {
	// Repo : Delete Clothes Used By Id
	rows, err := s.clothesUsedRepo.DeleteClothesUsedByClothesId(id)
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// Task Scheduler
func (s *clothesUsedService) SchedulerGetUsedClothesReadyToWash(days int) ([]models.SchedulerUsedClothesReadyToWash, error) {
	// Repo : Get Used Clothes Ready To Wash
	rows, err := s.clothesUsedRepo.SchedulerFindUsedClothesReadyToWash(days)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
