package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserWeather Interface
type UserWeatherRepository interface {
	Create(weather *models.UserWeather, userID uuid.UUID) error
}

// UserWeather Struct
type userWeatherRepository struct {
	db *gorm.DB
}

// UserWeather Constructor
func NewUserWeatherRepository(db *gorm.DB) UserWeatherRepository {
	return &userWeatherRepository{db: db}
}

func (r *userWeatherRepository) Create(weather *models.UserWeather, userID uuid.UUID) error {
	weather.ID = uuid.New()
	weather.CreatedAt = time.Now()
	weather.CreatedBy = userID

	// Query
	return r.db.Create(weather).Error
}
