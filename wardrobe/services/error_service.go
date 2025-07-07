package services

import (
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/utils"
)

// Error Interface
type ErrorService interface {
	GetAllError(pagination utils.Pagination) ([]models.ErrorAudit, int64, error)

	// For Scheduler
	SchedulerGetAllErrorAudit() ([]models.ErrorAudit, error)
}

// Error Struct
type errorService struct {
	errorRepo repositories.ErrorRepository
}

// Error Constructor
func NewErrorService(errorRepo repositories.ErrorRepository) ErrorService {
	return &errorService{
		errorRepo: errorRepo,
	}
}

func (s *errorService) GetAllError(pagination utils.Pagination) ([]models.ErrorAudit, int64, error) {
	return s.errorRepo.FindAllError(pagination)
}

// For Scheduler
func (s *errorService) SchedulerGetAllErrorAudit() ([]models.ErrorAudit, error) {
	return s.errorRepo.FindAllErrorAudit()
}
