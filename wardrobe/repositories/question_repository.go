package repositories

import (
	"errors"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Question Interface
type QuestionRepository interface {
	CreateQuestion(question *models.Question) error
	FindAllQuestion() ([]models.Question, error)
	FindUnansweredQuestion() ([]models.UnansweredQuestion, error)
}

// Question Struct
type questionRepository struct {
	db *gorm.DB
}

// Question Constructor
func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) CreateQuestion(question *models.Question) error {
	// Default
	question.ID = uuid.New()
	question.IsShow = false
	question.Answer = nil

	// Query
	return r.db.Create(question).Error
}

func (r *questionRepository) FindAllQuestion() ([]models.Question, error) {
	// Model
	var questions []models.Question

	// Query
	if err := r.db.Find(&questions).Error; err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return questions, nil
}

func (r *questionRepository) FindUnansweredQuestion() ([]models.UnansweredQuestion, error) {
	// Model
	var question []models.UnansweredQuestion

	// Query
	result := r.db.Table("questions").
		Select("question, created_at").
		Where("answer is null").
		Order("created_at DESC").
		Find(&question)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(question) == 0 {
		return nil, errors.New("unanswered question not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return question, nil
}
