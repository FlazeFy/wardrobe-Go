package repositories

import (
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Feedback Interface
type FeedbackRepository interface {
	CreateFeedback(feedback *models.Feedback, userID uuid.UUID) error
	FindAllFeedback() ([]models.Feedback, error)
	HardDeleteFeedbackByID(ID uuid.UUID) error

	// For Feedback
	DeleteAll() error
}

// Feedback Struct
type feedbackRepository struct {
	db *gorm.DB
}

// Feedback Constructor
func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepository{db: db}
}

func (r *feedbackRepository) FindAllFeedback() ([]models.Feedback, error) {
	// Model
	var feedbacks []models.Feedback

	// Query
	if err := r.db.Preload("User").Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	if len(feedbacks) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return feedbacks, nil
}

func (r *feedbackRepository) CreateFeedback(feedback *models.Feedback, userID uuid.UUID) error {
	// Default
	feedback.ID = uuid.New()
	feedback.CreatedAt = time.Now()
	feedback.CreatedBy = userID

	// Query
	return r.db.Create(feedback).Error
}

func (r *feedbackRepository) HardDeleteFeedbackByID(ID uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("id = ?", ID).Delete(&models.Feedback{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// For Seeder
func (r *feedbackRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.Feedback{}).Error
}
