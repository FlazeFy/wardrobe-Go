package repositories

import (
	"errors"
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Clothes Used Interface
type ClothesUsedRepository interface {
	CreateClothesUsed(clothesUsed *models.ClothesUsed, userID uuid.UUID) error
	FindClothesUsedHistory(userID uuid.UUID, clothesID uuid.UUID, order string) ([]models.ClothesUsedHistory, error)
	HardDeleteClothesUsedByID(ID, createdBy uuid.UUID) error
	HardDeleteClothesUsedByClothesID(clothesID, createdBy uuid.UUID) error
	DeleteClothesUsedByClothesId(id uuid.UUID) (int64, error)

	// Task Scheduler
	SchedulerFindUsedClothesReadyToWash(days int) ([]models.SchedulerUsedClothesReadyToWash, error)
}

// Clothes Used Struct
type clothesUsedRepository struct {
	db *gorm.DB
}

// Clothes Used Constructor
func NewClothesUsedRepository(db *gorm.DB) ClothesUsedRepository {
	return &clothesUsedRepository{db: db}
}

func (r *clothesUsedRepository) CreateClothesUsed(clothesUsed *models.ClothesUsed, userID uuid.UUID) error {
	// Default
	clothesUsed.ID = uuid.New()
	clothesUsed.CreatedAt = time.Now()
	clothesUsed.CreatedBy = userID

	// Query
	return r.db.Create(clothesUsed).Error
}

func (r *clothesUsedRepository) FindClothesUsedHistory(userID uuid.UUID, clothesID uuid.UUID, order string) ([]models.ClothesUsedHistory, error) {
	// Model
	var clothes []models.ClothesUsedHistory

	// Ordering Prep
	is_desc := true
	if order == "asc" {
		is_desc = false
	}

	// Query
	query := r.db.Table("clothes_useds").
		Select("clothes_useds.id, clothes_name, clothes_type, clothes_note, used_context, clothes.created_at").
		Joins("JOIN clothes ON clothes.id = clothes_useds.clothes_id").
		Where("clothes_useds.created_by = ?", userID)

	if clothesID != uuid.Nil {
		query = query.Where("clothes_id = ?", clothesID)
	}

	query = query.Order(clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "clothes_useds.created_at"}, Desc: is_desc},
			{Column: clause.Column{Name: "clothes_name"}, Desc: is_desc},
		},
	})

	result := query.Find(&clothes)

	// Response
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}

func (r *clothesUsedRepository) HardDeleteClothesUsedByID(ID, createdBy uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("id = ?", ID).Where("created_by = ?", createdBy).Delete(&models.ClothesUsed{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *clothesUsedRepository) HardDeleteClothesUsedByClothesID(clothesID, createdBy uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("clothes_id = ? AND created_by = ?", clothesID, createdBy).Delete(&models.ClothesUsed{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *clothesUsedRepository) DeleteClothesUsedByClothesId(id uuid.UUID) (int64, error) {
	// Model
	var clothes models.ClothesUsed

	// Query
	result := r.db.Unscoped().Where("clothes_id", id).Delete(&clothes)

	// Response
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// Task Scheduler
func (r *clothesUsedRepository) SchedulerFindUsedClothesReadyToWash(days int) ([]models.SchedulerUsedClothesReadyToWash, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// Model
	var clothes []models.SchedulerUsedClothesReadyToWash

	// Subquery: latest usage per clothes_id
	latestUsage := r.db.Table("clothes_useds").
		Select("clothes_id, MAX(created_at) AS latest_used_at").
		Group("clothes_id")

	// Main Query
	result := r.db.Table("clothes_useds").
		Select(`clothes.clothes_name, clothes.clothes_type, clothes.clothes_made_from,
			clothes_useds.used_context, clothes.is_faded, clothes.is_scheduled,
			clothes_useds.created_at, users.username, users.telegram_is_valid,
			users.telegram_user_id`).
		Joins("JOIN (?) AS latest_usage ON clothes_useds.clothes_id = latest_usage.clothes_id AND clothes_useds.created_at = latest_usage.latest_used_at", latestUsage).
		Joins("JOIN clothes ON clothes.id = clothes_useds.clothes_id").
		Joins("JOIN users ON users.id = clothes_useds.created_by").
		Where("clothes_useds.created_at < ?", cutoffDate).
		Where("NOT EXISTS (?)", r.db.
			Table("washes").
			Select("1").
			Where("washes.clothes_id = clothes_useds.clothes_id").
			Where("washes.finished_at BETWEEN clothes_useds.created_at AND NOW()")).
		Order("users.username ASC").
		Order("clothes_useds.created_at DESC").
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
