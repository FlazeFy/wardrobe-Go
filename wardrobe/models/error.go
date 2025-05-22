package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ErrorContext struct {
	DB *gorm.DB
}

func NewErrorContext(db *gorm.DB) *ErrorContext {
	return &ErrorContext{DB: db}
}

type (
	Error struct {
		ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		Message    string    `json:"message" gorm:"type:text;not null"`
		StackTrace string    `json:"stack_trace" gorm:"type:text;not null"`
		File       string    `json:"file" gorm:"type:varchar(255);not null"`
		Line       uint      `json:"line" gorm:"not null"`
		IsFixed    bool      `json:"is_fixed" gorm:"not null"`
		CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	ErrorAudit struct {
		Message   string `json:"message"`
		CreatedAt string `json:"created_at"`
		Total     int    `json:"total"`
	}
)

func (c *ErrorContext) GetAllErrorAudit() ([]ErrorAudit, error) {
	// Model
	var errors_list []ErrorAudit

	// Query
	result := c.DB.Table("errors").
		Select("message, string_agg(created_at::text, ', ') as created_at, COUNT(1) as total").
		Group("message").
		Order(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{Column: clause.Column{Name: "total"}, Desc: true},
				{Column: clause.Column{Name: "message"}, Desc: false},
				{Column: clause.Column{Name: "created_at"}, Desc: false},
			},
		}).Find(&errors_list)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(errors_list) == 0 {
		return nil, errors.New("error not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return errors_list, nil
}
