package services

import (
	"wardrobe/models/others"
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// Stats Interface
type StatsService interface {
	GetMostUsedContext(tableName, targetCol string, userId uuid.UUID) ([]others.StatsContextTotal, error)
}

// Stats Struct
type statsService struct {
	statsRepo repositories.StatsRepository
}

// Stats Constructor
func NewStatsService(statsRepo repositories.StatsRepository) StatsService {
	return &statsService{
		statsRepo: statsRepo,
	}
}

func (s *statsService) GetMostUsedContext(tableName, targetCol string, userId uuid.UUID) ([]others.StatsContextTotal, error) {
	return s.statsRepo.FindMostUsedContext(tableName, targetCol, userId)
}
