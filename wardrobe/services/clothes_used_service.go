package services

import (
	"wardrobe/models"
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// Clothes Used Interface
type ClothesUsedService interface {
	CreateClothesUsed(clothesUsed *models.ClothesUsed, userID uuid.UUID) error
	GetClothesUsedHistory(userID uuid.UUID, clothesID uuid.UUID, order string) ([]models.ClothesUsedHistory, error)
	HardDeleteClothesUsedByID(ID, createdBy uuid.UUID) error
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

func (r *clothesUsedService) GetClothesUsedHistory(userID uuid.UUID, clothesID uuid.UUID, order string) ([]models.ClothesUsedHistory, error) {
	return r.clothesUsedRepo.FindClothesUsedHistory(userID, clothesID, order)
}

func (r *clothesUsedService) CreateClothesUsed(clothesUsed *models.ClothesUsed, userID uuid.UUID) error {
	return r.clothesUsedRepo.CreateClothesUsed(clothesUsed, userID)
}

func (r *clothesUsedService) HardDeleteClothesUsedByID(ID, createdBy uuid.UUID) error {
	return r.clothesUsedRepo.HardDeleteClothesUsedByID(ID, createdBy)
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
