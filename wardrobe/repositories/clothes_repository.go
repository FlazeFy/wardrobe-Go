package repositories

import (
	"errors"
	"time"
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Clothes Interface
type ClothesRepository interface {
	CreateClothes(clothes *models.Clothes, userID uuid.UUID) (*models.Clothes, error)
	FindClothesById(ID uuid.UUID) (*models.Clothes, error)
	CheckClothesByName(clothesName string, userID uuid.UUID) (bool, error)
	FindAllClothesHeader(pagination utils.Pagination, category, order string, userID uuid.UUID) ([]models.ClothesHeader, int64, error)
	FindAllClothesDetail(pagination utils.Pagination, category, order string, userID uuid.UUID) ([]models.Clothes, int64, error)
	FindClothesShortInfoById(id uuid.UUID) (*models.ClothesShortInfo, error)
	FindClothesLastCreated(ctx string, userID uuid.UUID) (*models.ClothesLastCreated, error)
	FindClothesLastDeleted(ctx string, userID uuid.UUID) (*models.ClothesLastDeleted, error)
	FindDeletedClothes(userID uuid.UUID) ([]models.ClothesDeleted, error)
	FindClothesPlanDestroy(days int) ([]models.ClothesPlanDestroy, error)
	UpdateClothesById(clothes *models.Clothes, ID uuid.UUID) error
	HardDeleteClothesById(id, createdBy uuid.UUID) error
	HardDeleteClothesById2(id uuid.UUID) (int64, error)

	// Task Scheduler
	SchedulerFindUnusedClothes(days int) ([]models.SchedulerClothesUnused, error)
	SchedulerFindUnironedClothes() ([]models.SchedulerClothesUnironed, error)
}

// Clothes Struct
type clothesRepository struct {
	db *gorm.DB
}

// Clothes Constructor
func NewClothesRepository(db *gorm.DB) ClothesRepository {
	return &clothesRepository{db: db}
}

func (r *clothesRepository) FindClothesById(ID uuid.UUID) (*models.Clothes, error) {
	// Model
	var clothes models.Clothes

	// Query
	result := r.db.Unscoped().Where("id = ?", ID).First(&clothes)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}

	return &clothes, nil
}

func (r *clothesRepository) CheckClothesByName(clothesName string, userID uuid.UUID) (bool, error) {
	// Model
	var clothes models.Clothes

	// Query
	result := r.db.Unscoped().Where("LOWER(clothes_name) = LOWER(?) AND created_by = ?", clothesName, userID).First(&clothes)

	// Response
	if result.Error != nil {
		return true, result.Error
	}
	if clothes.ID == uuid.Nil {
		return false, nil
	}

	return true, nil
}

func (r *clothesRepository) FindAllClothesHeader(pagination utils.Pagination, category, order string, userID uuid.UUID) ([]models.ClothesHeader, int64, error) {
	// Model
	var total int64
	var clothes []models.ClothesHeader

	// Ordering Prep
	is_desc := true
	if order == "asc" {
		is_desc = false
	}

	// Pagination Count
	offset := (pagination.Page - 1) * pagination.Limit
	countQuery := r.db.Model(&models.ClothesHeader{}).
		Where("created_by = ?", userID).
		Where("deleted_at IS NULL")
	if category != "all" {
		countQuery = countQuery.Where("clothes_category = ?", category)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Model
	query := r.db.Table("clothes").
		Select("id,clothes_name,clothes_image,clothes_size,clothes_gender,clothes_color,clothes_category,clothes_type,clothes_qty,is_faded,has_washed,has_ironed,is_favorite,is_scheduled").
		Where("created_by = ?", userID).
		Where("deleted_at is null").
		Limit(pagination.Limit).
		Offset(offset)

	if category != "all" {
		query = query.Where("clothes_category = ?", category)
	}

	query = query.Order(clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "is_favorite"}, Desc: true},
			{Column: clause.Column{Name: "clothes_name"}, Desc: is_desc},
			{Column: clause.Column{Name: "created_at"}, Desc: is_desc},
		},
	})

	result := query.Find(&clothes)

	// Response
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return clothes, total, nil
}

func (r *clothesRepository) CreateClothes(clothes *models.Clothes, userID uuid.UUID) (*models.Clothes, error) {
	// Default
	clothes.ID = uuid.New()
	clothes.CreatedAt = time.Now()
	clothes.CreatedBy = userID
	clothes.UpdatedAt = nil
	clothes.DeletedAt = nil

	// Query
	if err := r.db.Create(clothes).Error; err != nil {
		return nil, err
	}

	return clothes, nil
}

func (r *clothesRepository) UpdateClothesById(clothes *models.Clothes, ID uuid.UUID) error {
	now := time.Now()
	clothes.ID = ID
	clothes.UpdatedAt = &now

	if err := r.db.Save(clothes).Error; err != nil {
		return err
	}

	return nil
}

func (r *clothesRepository) FindAllClothesDetail(pagination utils.Pagination, category, order string, userID uuid.UUID) ([]models.Clothes, int64, error) {
	// Model
	var total int64
	var clothes []models.Clothes

	// Ordering Prep
	is_desc := true
	if order == "asc" {
		is_desc = false
	}

	// Pagination Count
	offset := (pagination.Page - 1) * pagination.Limit
	countQuery := r.db.Model(&models.ClothesHeader{}).
		Where("created_by = ?", userID).
		Where("deleted_at IS NULL")
	if category != "all" {
		countQuery = countQuery.Where("clothes_category = ?", category)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Model
	query := r.db.Table("clothes").
		Where("created_by = ?", userID).
		Where("deleted_at is null").
		Limit(pagination.Limit).
		Offset(offset)

	if category != "all" {
		query = query.Where("clothes_category = ?", category)
	}

	query = query.Order(clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "is_favorite"}, Desc: true},
			{Column: clause.Column{Name: "clothes_name"}, Desc: is_desc},
			{Column: clause.Column{Name: "created_at"}, Desc: is_desc},
		},
	})

	result := query.Find(&clothes)

	// Response
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return clothes, total, nil
}

func (r *clothesRepository) FindClothesShortInfoById(id uuid.UUID) (*models.ClothesShortInfo, error) {
	// Model
	var clothes models.ClothesShortInfo

	// Query
	result := r.db.Table("clothes").
		Select("clothes_name,clothes_type,clothes_category").
		Where("id = ?", id).
		First(&clothes)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}

	return &clothes, nil
}

func (r *clothesRepository) FindClothesLastCreated(ctx string, userID uuid.UUID) (*models.ClothesLastCreated, error) {
	// Model
	var clothes models.ClothesLastCreated

	// Query
	result := r.db.Table("clothes").
		Select("clothes_name, "+ctx).
		Where("created_by = ?", userID).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx},
			Desc:   true,
		}).First(&clothes)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}

	return &clothes, nil
}

func (r *clothesRepository) FindClothesLastDeleted(ctx string, userID uuid.UUID) (*models.ClothesLastDeleted, error) {
	// Model
	var clothes models.ClothesLastDeleted

	// Query
	result := r.db.Table("clothes").
		Select("clothes_name, "+ctx).
		Where("created_by = ?", userID).
		Where("deleted_at IS NOT NULL").
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx},
			Desc:   true,
		}).First(&clothes)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}

	return &clothes, nil
}

func (r *clothesRepository) FindDeletedClothes(userID uuid.UUID) ([]models.ClothesDeleted, error) {
	// Model
	var clothes []models.ClothesDeleted

	// Query
	result := r.db.Table("clothes").
		Select("id, clothes_name, clothes_image, clothes_size, clothes_gender, clothes_color, clothes_category, clothes_type, clothes_qty, deleted_at").
		Where("created_by = ?", userID).
		Where("deleted_at IS NOT NULL").
		Order("deleted_at DESC").
		Find(&clothes)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}

func (r *clothesRepository) FindClothesPlanDestroy(days int) ([]models.ClothesPlanDestroy, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// Model
	var clothes []models.ClothesPlanDestroy

	// Query
	result := r.db.Table("clothes").
		Select("clothes.id,clothes_name,username,telegram_user_id,telegram_is_valid").
		Joins("JOIN users ON users.id = clothes.created_by").
		Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffDate).
		Order("username ASC").
		Find(&clothes)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(clothes) == 0 {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}

func (r *clothesRepository) HardDeleteClothesById(id, createdBy uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("id = ? AND created_by = ? AND deleted_at IS NOT NULL", id, createdBy).Delete(&models.Clothes{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// fix this to use user id
func (r *clothesRepository) HardDeleteClothesById2(id uuid.UUID) (int64, error) {
	// Model
	var clothes models.Clothes

	// Query
	result := r.db.Unscoped().Where("id", id).Delete(&clothes)

	// Response
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// Command Scheduler
func (r *clothesRepository) SchedulerFindUnusedClothes(days int) ([]models.SchedulerClothesUnused, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// Model
	var clothes []models.SchedulerClothesUnused

	// Query
	result := r.db.Table("clothes").
		Select(`
			clothes.clothes_name,clothes.clothes_type,
			COALESCE(MAX(clothes_useds.created_at), clothes.created_at) AS last_used,
			COUNT(clothes_useds.id) AS total_used,
			users.username,users.telegram_user_id,users.telegram_is_valid
		`).
		Joins("JOIN users ON users.id = clothes.created_by").
		Joins("LEFT JOIN clothes_useds ON clothes.id = clothes_useds.clothes_id").
		Group("clothes.id,clothes.clothes_name,clothes.clothes_type,clothes.created_at,users.username,users.telegram_user_id,users.telegram_is_valid").
		Having("COALESCE(MAX(clothes_useds.created_at), clothes.created_at) < ?", cutoffDate).
		Order("username ASC").
		Find(&clothes)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(clothes) == 0 {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}

func (r *clothesRepository) SchedulerFindUnironedClothes() ([]models.SchedulerClothesUnironed, error) {
	ironable_clothes_made_from := []string{"cotton", "linen", "silk", "rayon"}
	ironable_clothes_type := []string{"pants", "shirt", "jacket", "shorts", "skirt", "dress", "blouse", "sweater", "hoodie", "tie", "coat", "vest", "t-shirt", "jeans", "leggings", "cardigan"}

	// Model
	var clothes []models.SchedulerClothesUnironed

	// Query
	result := r.db.Table("clothes").
		Select("clothes_name,clothes_made_from,has_washed,is_favorite,is_scheduled,username,telegram_user_id,telegram_is_valid").
		Joins("JOIN users ON users.id = clothes.created_by").
		Where("has_ironed = ?", false).
		Where("clothes_made_from IN (?)", ironable_clothes_made_from).
		Where("clothes_type IN (?)", ironable_clothes_type).
		Order("username ASC").
		Find(&clothes)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(clothes) == 0 {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}
