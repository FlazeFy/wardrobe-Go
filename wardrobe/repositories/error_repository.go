package repositories

import (
	"errors"
	"wardrobe/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Error Interface
type ErrorRepository interface {
	GetAllErrorAudit() ([]models.ErrorAudit, error)
}

// Error Struct
type errorRepository struct {
	db *gorm.DB
}

// Error Constructor
func NewErrorRepository(db *gorm.DB) ErrorRepository {
	return &errorRepository{db: db}
}

func (r *errorRepository) GetAllErrorAudit() ([]models.ErrorAudit, error) {
	// Model
	var errors_list []models.ErrorAudit

	// Query
	result := r.db.Table("errors").
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
