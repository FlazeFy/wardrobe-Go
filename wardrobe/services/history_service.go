package services

import (
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/utils"

	"github.com/google/uuid"
)

// History Interface
type HistoryService interface {
	GetAllHistory(pagination utils.Pagination, userID *uuid.UUID) ([]models.GetHistory, int64, error)
	HardDeleteHistoryByID(ID, createdBy uuid.UUID) error

	// Task Scheduler
	DeleteHistoryForLastNDays(days int) (int64, error)
}

// History Struct
type historyService struct {
	historyRepo repositories.HistoryRepository
}

// History Constructor
func NewHistoryService(historyRepo repositories.HistoryRepository) HistoryService {
	return &historyService{
		historyRepo: historyRepo,
	}
}

func (r *historyService) GetAllHistory(pagination utils.Pagination, userID *uuid.UUID) ([]models.GetHistory, int64, error) {
	return r.historyRepo.FindAllHistory(pagination, userID)
}

func (r *historyService) HardDeleteHistoryByID(ID, createdBy uuid.UUID) error {
	return r.historyRepo.HardDeleteHistoryByID(ID, createdBy)
}

// For Task Scheduler
func (r *historyService) DeleteHistoryForLastNDays(days int) (int64, error) {
	rows, err := r.historyRepo.DeleteHistoryForLastNDays(days)
	if err != nil {
		return 0, err
	}

	return rows, nil
}
