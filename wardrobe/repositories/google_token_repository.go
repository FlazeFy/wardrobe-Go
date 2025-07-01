package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Google Token Interface
type GoogleTokenRepository interface {
	CreateGoogleToken(googleToken *models.GoogleToken, userID uuid.UUID) error
}

// Google Token Struct
type googleTokenRepository struct {
	db *gorm.DB
}

// Google Token Constructor
func NewGoogleTokenRepository(db *gorm.DB) GoogleTokenRepository {
	return &googleTokenRepository{db: db}
}

func (r *googleTokenRepository) CreateGoogleToken(googleToken *models.GoogleToken, userID uuid.UUID) error {
	// Default
	googleToken.ID = uuid.New()
	googleToken.CreatedAt = time.Now()
	googleToken.CreatedBy = userID

	// Query
	return r.db.Create(googleToken).Error
}
