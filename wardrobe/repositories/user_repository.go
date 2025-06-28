package repositories

import (
	"errors"
	"time"
	"wardrobe/models"
	"wardrobe/models/others"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Interface
type UserRepository interface {
	FindByUsernameOrEmail(username, email string) (*models.User, error)
	FindUserContactByID(id uuid.UUID) (*models.UserContact, error)
	FindByEmail(email string) (*models.User, error)
	FindById(id string) (*others.MyProfile, error)
	CreateUser(user *models.User) (*models.User, error)

	// For Task Scheduler
	FindUserReadyFetchWeather() ([]models.UserReadyFetchWeather, error)

	// For Seeder
	DeleteAll() error
}

// User Struct
type userRepository struct {
	db *gorm.DB
}

// User Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsernameOrEmail(username, email string) (*models.User, error) {
	// Models
	var user models.User

	// Query
	err := r.db.Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	// Models
	var user models.User

	// Query
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindById(id string) (*others.MyProfile, error) {
	// Models
	var user others.MyProfile

	// Query
	err := r.db.Table("users").
		Select("username, email, telegram_is_valid, telegram_user_id, created_at").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.TelegramIsValid = false

	// Query
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// For Task Scheduler
func (r *userRepository) FindUserReadyFetchWeather() ([]models.UserReadyFetchWeather, error) {
	// Model
	var users []models.UserReadyFetchWeather

	// Query
	result := r.db.Table("users").
		Select(`
			user_tracks.track_lat,user_tracks.track_long,user_tracks.created_at,users.id as user_id,
			users.username,users.telegram_user_id,users.telegram_is_valid
		`).
		Joins(`
			JOIN (
				SELECT DISTINCT ON (created_by) *
				FROM user_tracks
				ORDER BY created_by, created_at DESC
			) AS user_tracks ON user_tracks.created_by = users.id
		`).
		Where("user_tracks.track_lat IS NOT NULL AND user_tracks.track_long IS NOT NULL").
		Order("users.username ASC").
		Find(&users)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}
	if len(users) == 0 {
		return nil, errors.New("no user track found")
	}

	return users, nil
}

func (r *userRepository) FindUserContactByID(id uuid.UUID) (*models.UserContact, error) {
	// Model
	var contact models.UserContact

	// Query
	result := r.db.Table("users").
		Select("username, email, telegram_user_id, telegram_is_valid").
		Where("id = ?", id).
		First(&contact)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}

	return &contact, nil
}

// For Seeder
func (r *userRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.Admin{}).Error
}
