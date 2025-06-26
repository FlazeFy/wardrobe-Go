package services

import (
	"wardrobe/models"
	"wardrobe/repositories"
)

// Question Interface
type QuestionService interface {
	GetUnansweredQuestion() ([]models.UnansweredQuestion, error)
}

// Question Struct
type questionService struct {
	questionRepo repositories.QuestionRepository
}

// Question Constructor
func NewQuestionService(questionRepo repositories.QuestionRepository) QuestionService {
	return &questionService{
		questionRepo: questionRepo,
	}
}

// For Task Scheduler
func (r *questionService) GetUnansweredQuestion() ([]models.UnansweredQuestion, error) {
	rows, err := r.questionRepo.FindUnansweredQuestion()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
