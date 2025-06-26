package services

import (
	"wardrobe/models"
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// History Interface
type HistoryService interface {
	GetAllHistory() ([]models.History, error)
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

func (r *historyService) GetAllHistory() ([]models.History, error) {
	return r.historyRepo.FindAllHistory()
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
