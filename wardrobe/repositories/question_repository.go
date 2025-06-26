package repositories

import (
	"errors"
	"wardrobe/models"

	"gorm.io/gorm"
)

// Question Interface
type QuestionRepository interface {
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

func (r *questionRepository) FindUnansweredQuestion() ([]models.UnansweredQuestion, error) {
	// Model
	var question []models.UnansweredQuestion

	// Query
	result := r.db.Table("questions").
		Select("question, created_at").
		Where("answer is null").
		Order("created_at DESC").
		Find(&question)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(question) == 0 {
		return nil, errors.New("unanswered question not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return question, nil
}
