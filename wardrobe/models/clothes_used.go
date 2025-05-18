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
)

func (c *ClothesUsedContext) GetClothesUsedHistory(userID uuid.UUID, clothesID uuid.UUID, order string) ([]ClothesUsedHistory, error) {
	var clothes []ClothesUsedHistory

	is_desc := true
	if order == "asc" {
		is_desc = false
	}

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

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(clothes) == 0 {
		return nil, errors.New("clothes used history not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}
