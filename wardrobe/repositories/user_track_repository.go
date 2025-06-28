package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Track Interface
type UserTrackRepository interface {
	CreateUserTrack(track *models.UserTrack, userID uuid.UUID) error

	// For Seeder
	DeleteAll() error
}

// User Track Struct
type userTrackRepository struct {
	db *gorm.DB
}

// User Track Constructor
func NewUserTrackRepository(db *gorm.DB) UserTrackRepository {
	return &userTrackRepository{db: db}
}

func (r *userTrackRepository) CreateUserTrack(track *models.UserTrack, userID uuid.UUID) error {
	track.ID = uuid.New()
	track.CreatedAt = time.Now()
	track.CreatedBy = userID

	// Query
	return r.db.Create(track).Error
}

// For Seeder
func (r *userTrackRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.UserTrack{}).Error
}
