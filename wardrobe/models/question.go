package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionContext struct {
	DB *gorm.DB
}

func NewQuestionContext(db *gorm.DB) *QuestionContext {
	return &QuestionContext{DB: db}
}

type (
	Question struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Question  string    `json:"question" gorm:"type:varchar(500);not null"`
		Answer    *string   `json:"answer" gorm:"type:varchar(500);null"`
		IsShow    bool      `json:"is_show" gorm:"type:boolean;not null"`
		CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	UnansweredQuestion struct {
		Question  string    `json:"question" gorm:"type:varchar(500);not null"`
		CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
)

func (c *QuestionContext) GetUnansweredQuestion() ([]UnansweredQuestion, error) {
	// Model
	var question []UnansweredQuestion

	// Query
	result := c.DB.Table("questions").
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
