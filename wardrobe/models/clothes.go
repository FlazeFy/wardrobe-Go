package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClothesContext struct {
	DB *gorm.DB
}

func NewClothesContext(db *gorm.DB) *ClothesContext {
	return &ClothesContext{DB: db}
}

type (
	Clothes struct {
		ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesName  string     `json:"clothes_name" gorm:"type:varchar(36);not null"`
		ClothesDesc  *string    `json:"clothes_desc" gorm:"type:varchar(500);null"`
		ClothesMerk  *string    `json:"clothes_merk" gorm:"type:varchar(75);null"`
		ClothesColor string     `json:"clothes_color" gorm:"type:varchar(36);not null"`
		ClothesPrice *int       `json:"clothes_price" gorm:"type:int;null"`
		ClothesBuyAt *time.Time `json:"clothes_buy_at" gorm:"null"`
		ClothesQty   int        `json:"clothes_qty" gorm:"type:int;not null"`
		ClothesImage *string    `json:"clothes_image" gorm:"type:varchar(1000);null"`
		IsFaded      bool       `json:"is_faded" gorm:"type:boolean;not null"`
		HasWashed    bool       `json:"has_washed" gorm:"type:boolean;not null"`
		HasIroned    bool       `json:"has_ironed" gorm:"type:boolean;not null"`
		IsFavorite   bool       `json:"is_favorite" gorm:"type:boolean;not null"`
		IsScheduled  bool       `json:"is_scheduled" gorm:"type:boolean;not null"`
		CreatedAt    time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
		UpdatedAt    *time.Time `json:"updated_at" gorm:"type:timestamp;null"`
		DeletedAt    *time.Time `json:"deleted_at" gorm:"type:timestamp;null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Dictionary
		ClothesMadeFrom    string     `json:"clothes_made_from" gorm:"type:varchar(36);not null"`
		DictionaryMadeFrom Dictionary `json:"-" gorm:"foreignKey:ClothesMadeFrom;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesType        string     `json:"clothes_type" gorm:"type:varchar(36);not null"`
		DictionaryType     Dictionary `json:"-" gorm:"foreignKey:ClothesType;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesCategory    string     `json:"clothes_category" gorm:"type:varchar(36);not null"`
		DictionaryCategory Dictionary `json:"-" gorm:"foreignKey:ClothesCategory;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesSize        string     `json:"clothes_size" gorm:"type:varchar(3);not null"`
		DictionarySize     Dictionary `json:"-" gorm:"foreignKey:ClothesSize;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesGender      string     `json:"clothes_gender" gorm:"type:varchar(6);not null"`
		DictionaryGender   Dictionary `json:"-" gorm:"foreignKey:ClothesGender;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	ClothesShortInfo struct {
		ClothesName string `json:"clothes_name" gorm:"type:varchar(36);not null"`
		// FK - Dictionary
		ClothesType        string     `json:"clothes_type" gorm:"type:varchar(36);not null"`
		DictionaryType     Dictionary `json:"-" gorm:"foreignKey:ClothesType;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesCategory    string     `json:"clothes_category" gorm:"type:varchar(36);not null"`
		DictionaryCategory Dictionary `json:"-" gorm:"foreignKey:ClothesCategory;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	ClothesLastCreated struct {
		ClothesName string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	ClothesLastDeleted struct {
		ClothesName string     `json:"clothes_name" gorm:"type:varchar(36);not null"`
		DeletedAt   *time.Time `json:"deleted_at" gorm:"type:timestamp;null"`
	}
)

func (c *ClothesContext) GetClothesShortInfoById(id uuid.UUID) (*ClothesShortInfo, error) {
	var clothes ClothesShortInfo
	result := c.DB.Table("clothes").
		Select("clothes_name,clothes_type,clothes_category").
		Where("id = ?", id).
		First(&clothes)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &clothes, nil
}

func (c *ClothesContext) GetClothesLastCreated(ctx string, userID uuid.UUID) (*ClothesLastCreated, error) {
	var clothes ClothesLastCreated

	result := c.DB.Table("clothes").
		Select("clothes_name, "+ctx).
		Where("created_by = ?", userID).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx},
			Desc:   true,
		}).First(&clothes)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &clothes, nil
}

func (c *ClothesContext) GetClothesLastDeleted(ctx string, userID uuid.UUID) (*ClothesLastDeleted, error) {
	var clothes ClothesLastDeleted

	result := c.DB.Table("clothes").
		Select("clothes_name, "+ctx).
		Where("created_by = ?", userID).
		Where("deleted_at IS NOT NULL").
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx},
			Desc:   true,
		}).First(&clothes)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &clothes, nil
}
