package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClothesUsedContext struct {
	DB *gorm.DB
}

func NewClothesUsedContext(db *gorm.DB) *ClothesUsedContext {
	return &ClothesUsedContext{DB: db}
}

type (
	ClothesUsed struct {
		ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesNote *string   `json:"clothes_note" gorm:"type:varchar(500);null"`
		CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Clothes
		ClothesId uuid.UUID `json:"clothes_id" gorm:"not null"`
		Clothes   Clothes   `json:"-" gorm:"foreignKey:ClothesId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Used Context
		UsedContext string     `json:"used_context" gorm:"not null"`
		Dictionary  Dictionary `json:"-" gorm:"foreignKey:UsedContext;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	ClothesUsedHistory struct {
		ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesName string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		ClothesType string    `json:"clothes_type" gorm:"type:varchar(36);not null"`
		ClothesNote *string   `json:"clothes_note" gorm:"type:varchar(500);null"`
		CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - Used Context
		UsedContext string `json:"used_context" gorm:"not null"`
	}
	SchedulerUsedClothesReadyToWash struct {
		ClothesName     string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		ClothesType     string    `json:"clothes_type" gorm:"type:varchar(36);not null"`
		ClothesMadeFrom string    `json:"clothes_made_from" gorm:"type:varchar(36);not null"`
		IsFaded         bool      `json:"is_faded" gorm:"type:boolean;not null"`
		IsScheduled     bool      `json:"is_scheduled" gorm:"type:boolean;not null"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		// FK - Used Context
		UsedContext string `json:"used_context" gorm:"not null"`
	}
)

func (c *ClothesUsedContext) GetClothesUsedHistory(userID uuid.UUID, clothesID uuid.UUID, order string) ([]ClothesUsedHistory, error) {
	// Model
	var clothes []ClothesUsedHistory

	// Ordering Prep
	is_desc := true
	if order == "asc" {
		is_desc = false
	}

	// Query
	query := c.DB.Table("clothes_useds").
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
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(clothes) == 0 {
		return nil, errors.New("clothes used history not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}

func (c *ClothesUsedContext) SchedulerGetUsedClothesReadyToWash(days int) ([]SchedulerUsedClothesReadyToWash, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// Model
	var clothes []SchedulerUsedClothesReadyToWash

	// Subquery: latest usage per clothes_id
	latestUsage := c.DB.Table("clothes_useds").
		Select("clothes_id, MAX(created_at) AS latest_used_at").
		Group("clothes_id")

	// Main Query
	result := c.DB.Table("clothes_useds").
		Select(`clothes.clothes_name, clothes.clothes_type, clothes.clothes_made_from,
			clothes_useds.used_context, clothes.is_faded, clothes.is_scheduled,
			clothes_useds.created_at, users.username, users.telegram_is_valid,
			users.telegram_user_id`).
		Joins("JOIN (?) AS latest_usage ON clothes_useds.clothes_id = latest_usage.clothes_id AND clothes_useds.created_at = latest_usage.latest_used_at", latestUsage).
		Joins("JOIN clothes ON clothes.id = clothes_useds.clothes_id").
		Joins("JOIN users ON users.id = clothes_useds.created_by").
		Where("clothes_useds.created_at < ?", cutoffDate).
		Where("NOT EXISTS (?)", c.DB.
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
