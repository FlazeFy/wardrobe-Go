package repositories

import (
	"time"
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Error Interface
type ErrorRepository interface {
	FindAllError(pagination utils.Pagination) ([]models.ErrorAudit, int64, error)
	HardDeleteErrorByID(ID uuid.UUID) error

	// For Scheduler
	FindAllErrorAudit() ([]models.ErrorAudit, error)

	// For Seeder
	DeleteAll() error
	CreateError(errData *models.Error) error
}

// Error Struct
type errorRepository struct {
	db *gorm.DB
}

// Error Constructor
func NewErrorRepository(db *gorm.DB) ErrorRepository {
	return &errorRepository{db: db}
}

func (r *errorRepository) FindAllError(pagination utils.Pagination) ([]models.ErrorAudit, int64, error) {
	// Model
	var errorsList []models.ErrorAudit
	var total int64

	// Pagination Count
	offset := (pagination.Page - 1) * pagination.Limit
	countQuery := r.db.Model(&models.Error{}).
		Group("message")
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

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
		}).
		Limit(pagination.Limit).
		Offset(offset).
		Find(&errorsList)

	if len(errorsList) == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return errorsList, total, nil
}

func (r *errorRepository) FindAllErrorAudit() ([]models.ErrorAudit, error) {
	// Model
	var errorsList []models.ErrorAudit

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
		}).Find(&errorsList)

	if len(errorsList) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return errorsList, nil
}

func (r *errorRepository) CreateError(errData *models.Error) error {
	// Default
	errData.ID = uuid.New()
	errData.CreatedAt = time.Now()
	errData.IsFixed = false

	// Query
	return r.db.Create(errData).Error
}

func (r *errorRepository) HardDeleteErrorByID(ID uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("id = ?", ID).Delete(&models.Error{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// For Seeder
func (r *errorRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.Error{}).Error
}
