package services

import (
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// Schedule Interface
type ScheduleService interface {
	DeleteScheduleByClothesId(id uuid.UUID) (int64, error)
}

// Schedule Struct
type scheduleService struct {
	scheduleRepo repositories.ScheduleRepository
}

// Schedule Constructor
func NewScheduleService(scheduleRepo repositories.ScheduleRepository) ScheduleService {
	return &scheduleService{
		scheduleRepo: scheduleRepo,
	}
}

func (s *scheduleService) DeleteScheduleByClothesId(id uuid.UUID) (int64, error) {
	// Repo : Delete Schedule By Clothes Id
	rows, err := s.scheduleRepo.DeleteScheduleByClothesId(id)
	if err != nil {
		return 0, err
	}

	return rows, nil
}
