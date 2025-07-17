package repositories

import (
	"time"
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Feedback Interface
type FeedbackRepository interface {
	CreateFeedback(feedback *models.Feedback, userID uuid.UUID) error
	FindAllFeedback(pagination utils.Pagination) ([]models.Feedback, int64, error)
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

func (r *feedbackRepository) FindAllFeedback(pagination utils.Pagination) ([]models.Feedback, int64, error) {
	// Model
	var feedbacks []models.Feedback
	var total int64

	// Pagination Count
	offset := (pagination.Page - 1) * pagination.Limit
	countQuery := r.db.Model(&models.Feedback{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Query
	result := r.db.Preload("User").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&feedbacks)

	if len(feedbacks) == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return feedbacks, total, nil
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
