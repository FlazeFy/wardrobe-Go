package services

import (
	"wardrobe/models"
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// Clothes Interface
type ClothesService interface {
	// For Task Scheduler
	GetClothesPlanDestroy(days int) ([]models.ClothesPlanDestroy, error)
	SchedulerHardDeleteClothesById(id uuid.UUID) (int64, error)
	SchedulerGetUnusedClothes(days int) ([]models.SchedulerClothesUnused, error)
	SchedulerGetUnironedClothes() ([]models.SchedulerClothesUnironed, error)
}

// Clothes Struct
type clothesService struct {
	clothesRepo repositories.ClothesRepository
}

// Clothes Constructor
func NewClothesService(clothesRepo repositories.ClothesRepository) ClothesService {
	return &clothesService{
		clothesRepo: clothesRepo,
	}
}

// For Task Scheduler
func (s *clothesService) GetClothesPlanDestroy(days int) ([]models.ClothesPlanDestroy, error) {
	// Repo : Find Clothes Plan Destroy
	rows, err := s.clothesRepo.FindClothesPlanDestroy(days)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *clothesService) SchedulerHardDeleteClothesById(id uuid.UUID) (int64, error) {
	// Repo : Hard Delete Clothes By Id
	rows, err := s.clothesRepo.HardDeleteClothesById(id)
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *clothesService) SchedulerGetUnusedClothes(days int) ([]models.SchedulerClothesUnused, error) {
	// Repo : Scheduler Find Unused Clothes
	rows, err := s.clothesRepo.SchedulerFindUnusedClothes(days)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *clothesService) SchedulerGetUnironedClothes() ([]models.SchedulerClothesUnironed, error) {
	// Repo : Scheduler Find Unironed Clothes
	rows, err := s.clothesRepo.SchedulerFindUnironedClothes()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
