package services

import (
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/utils"
)

// User Interface
type UserService interface {
	GetAllUser(pagination utils.Pagination, order, username string) ([]models.UserAnalytic, int64, error)

	// Task Scheduler
	SchedulerGetUserReadyFetchWeather() ([]models.UserReadyFetchWeather, error)
}

// User Struct
type userService struct {
	userRepo repositories.UserRepository
}

// User Constructor
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAllUser(pagination utils.Pagination, order, username string) ([]models.UserAnalytic, int64, error) {
	return s.userRepo.FindAllUser(pagination, order, username)
}

func (r *userService) SchedulerGetUserReadyFetchWeather() ([]models.UserReadyFetchWeather, error) {
	// Repo : Find User Ready Fetch Weather
	rows, err := r.userRepo.FindUserReadyFetchWeather()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
