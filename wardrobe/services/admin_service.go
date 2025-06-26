package services

import (
	"wardrobe/models"
	"wardrobe/repositories"
)

// Admin Interface
type AdminService interface {
	GetAllAdminContact() ([]models.UserContact, error)
}

// Admin Struct
type adminService struct {
	adminRepo repositories.AdminRepository
}

// Admin Constructor
func NewAdminService(adminRepo repositories.AdminRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

func (s *adminService) GetAllAdminContact() ([]models.UserContact, error) {
	// Repo : Get All Admin Audit
	admin, err := s.adminRepo.FindAllAdminContact()
	if err != nil {
		return nil, err
	}

	return admin, nil
}
