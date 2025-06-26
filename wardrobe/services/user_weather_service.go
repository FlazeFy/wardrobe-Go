package services

import (
	"wardrobe/models"
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// User Weather Interface
type UserWeatherService interface {
	// Task Scheduler
	Create(weather *models.UserWeather, userID uuid.UUID) error
}

// User Weather Struct
type userWeatherService struct {
	userWeatherRepo repositories.UserWeatherRepository
}

// User Weather Constructor
func NewUserWeatherService(userWeatherRepo repositories.UserWeatherRepository) UserWeatherService {
	return &userWeatherService{
		userWeatherRepo: userWeatherRepo,
	}
}

func (s *userWeatherService) Create(weather *models.UserWeather, userID uuid.UUID) error {
	// Repo : Create User Weather
	if err := s.userWeatherRepo.Create(weather, userID); err != nil {
		return err
	}

	return nil
}
