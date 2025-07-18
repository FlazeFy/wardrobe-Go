package services

import (
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/utils"
)

// Question Interface
type QuestionService interface {
	CreateQuestion(question *models.Question) error
	GetAllQuestion(pagination utils.Pagination) ([]models.Question, int64, error)
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

func (r *questionService) GetAllQuestion(pagination utils.Pagination) ([]models.Question, int64, error) {
	return r.questionRepo.FindAllQuestion(pagination)
}

func (r *questionService) CreateQuestion(question *models.Question) error {
	return r.questionRepo.CreateQuestion(question)
}

// For Task Scheduler
func (r *questionService) GetUnansweredQuestion() ([]models.UnansweredQuestion, error) {
	rows, err := r.questionRepo.FindUnansweredQuestion()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
