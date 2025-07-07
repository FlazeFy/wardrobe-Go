package repositories

import (
	"errors"
	"wardrobe/models"
	"wardrobe/models/others"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Admin Interface
type AdminRepository interface {
	FindByEmail(email string) (*models.Admin, error)
	FindAllContact() ([]models.AdminContact, error)
	FindAllAdminContact() ([]models.UserContact, error)
	FindById(id uuid.UUID) (*others.MyProfile, error)

	// For Seeder
	Create(room *models.Admin) error
	DeleteAll() error
}

// Admin Struct
type adminRepository struct {
	db *gorm.DB
}

// Admin Constructor
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) FindAllContact() ([]models.AdminContact, error) {
	// Models
	var admin []models.AdminContact

	// Query
	err := r.db.Table("admins").
		Select("username, email, telegram_is_valid, telegram_user_id").
		Where("telegram_is_valid = ?", true).
		Where("telegram_user_id IS NOT NULL").
		Order("username ASC").
		Find(&admin).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return admin, err
}

func (r *adminRepository) FindByEmail(email string) (*models.Admin, error) {
	// Models
	var admin models.Admin

	// Query
	err := r.db.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *adminRepository) FindById(id uuid.UUID) (*others.MyProfile, error) {
	// Models
	var admin others.MyProfile

	// Query
	err := r.db.Table("admins").
		Select("username, email, telegram_is_valid, telegram_user_id, created_at").
		Where("id = ?", id).
		First(&admin).Error
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

// For Task Scheduler
func (r *adminRepository) FindAllAdminContact() ([]models.UserContact, error) {
	// Model
	var contact []models.UserContact

	// Query
	result := r.db.Table("admins").
		Select("username, email, telegram_user_id, telegram_is_valid").
		Find(&contact)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(contact) == 0 {
		return nil, errors.New("admin contact not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return contact, nil
}

// For Seeder
func (r *adminRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.Admin{}).Error
}
func (r *adminRepository) Create(admin *models.Admin) error {
	admin.ID = uuid.New()

	// Query
	return r.db.Create(admin).Error
}
