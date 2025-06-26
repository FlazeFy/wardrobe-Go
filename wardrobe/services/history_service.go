package services

import (
	"wardrobe/repositories"
)

// History Interface
type HistoryService interface {
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

// For Task Scheduler
func (r *historyService) DeleteHistoryForLastNDays(days int) (int64, error) {
	rows, err := r.historyRepo.DeleteHistoryForLastNDays(days)
	if err != nil {
		return 0, err
	}

	return rows, nil
}
