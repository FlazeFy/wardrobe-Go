package services

import (
	"wardrobe/models"
	"wardrobe/repositories"

	"github.com/google/uuid"
)

// Feedback Interface
type FeedbackService interface {
	CreateFeedback(feedback *models.Feedback, userID uuid.UUID) error
	HardDeleteFeedbackByID(ID uuid.UUID) error
	GetAllFeedback() ([]models.Feedback, error)
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

func (r *feedbackService) GetAllFeedback() ([]models.Feedback, error) {
	return r.feedbackRepo.FindAllFeedback()
}

func (r *feedbackService) CreateFeedback(feedback *models.Feedback, userID uuid.UUID) error {
	return r.feedbackRepo.CreateFeedback(feedback, userID)
}

func (r *feedbackService) HardDeleteFeedbackByID(ID uuid.UUID) error {
	return r.feedbackRepo.HardDeleteFeedbackByID(ID)
}
