package repositories

import (
	"errors"
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Question Interface
type QuestionRepository interface {
	CreateQuestion(question *models.Question) error
	FindAllQuestion(pagination utils.Pagination) ([]models.Question, int64, error)
	FindUnansweredQuestion() ([]models.UnansweredQuestion, error)

	// For Seeder
	DeleteAll() error
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

func (r *questionRepository) FindAllQuestion(pagination utils.Pagination) ([]models.Question, int64, error) {
	// Model
	var questions []models.Question
	var total int64

	// Pagination Count
	offset := (pagination.Page - 1) * pagination.Limit
	countQuery := r.db.Model(&models.Question{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Query
	result := r.db.Limit(pagination.Limit).
		Offset(offset).
		Find(&questions)

	if len(questions) == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return questions, total, nil
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

// For Seeder
func (r *questionRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.Question{}).Error
}
