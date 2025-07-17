package services

import (
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/utils"

	"github.com/google/uuid"
)

// Feedback Interface
type FeedbackService interface {
	CreateFeedback(feedback *models.Feedback, userID uuid.UUID) error
	HardDeleteFeedbackByID(ID uuid.UUID) error
	GetAllFeedback(pagination utils.Pagination) ([]models.Feedback, int64, error)
}

// Feedback Struct
type feedbackService struct {
	feedbackRepo repositories.FeedbackRepository
}

// Feedback Constructor
func NewFeedbackService(feedbackRepo repositories.FeedbackRepository) FeedbackService {
	return &feedbackService{
		feedbackRepo: feedbackRepo,
	}
}

func (r *feedbackService) GetAllFeedback(pagination utils.Pagination) ([]models.Feedback, int64, error) {
	return r.feedbackRepo.FindAllFeedback(pagination)
}

func (r *feedbackService) CreateFeedback(feedback *models.Feedback, userID uuid.UUID) error {
	return r.feedbackRepo.CreateFeedback(feedback, userID)
}

func (r *feedbackService) HardDeleteFeedbackByID(ID uuid.UUID) error {
	return r.feedbackRepo.HardDeleteFeedbackByID(ID)
}
