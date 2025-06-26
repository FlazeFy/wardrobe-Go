package services

import (
	"wardrobe/models"
	"wardrobe/repositories"
)

// Error Interface
type ErrorService interface {
	GetAllErrorAudit() ([]models.ErrorAudit, error)
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

func (s *errorService) GetAllErrorAudit() ([]models.ErrorAudit, error) {
	// Repo : Get All Error Audit
	errors_list, err := s.errorRepo.FindAllErrorAudit()
	if err != nil {
		return nil, err
	}

	return errors_list, nil
}
