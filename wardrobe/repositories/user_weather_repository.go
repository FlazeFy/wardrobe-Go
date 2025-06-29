package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Weather Interface
type UserWeatherRepository interface {
	CreateUserWeather(weather *models.UserWeather, userID uuid.UUID) error

	// For Seeder
	DeleteAll() error
}

// User Weather Struct
type userWeatherRepository struct {
	db *gorm.DB
}

// User Weather Constructor
func NewUserWeatherRepository(db *gorm.DB) UserWeatherRepository {
	return &userWeatherRepository{db: db}
}

func (r *userWeatherRepository) CreateUserWeather(weather *models.UserWeather, userID uuid.UUID) error {
	weather.ID = uuid.New()
	weather.CreatedAt = time.Now()
	weather.CreatedBy = userID

	// Query
	return r.db.Create(weather).Error
}

// For Seeder
func (r *userWeatherRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.UserWeather{}).Error
}
