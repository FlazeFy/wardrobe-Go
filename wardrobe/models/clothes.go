package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Clothes struct {
		ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesName  string     `json:"clothes_name" gorm:"type:varchar(75);not null" binding:"required,max=75"`
		ClothesDesc  *string    `json:"clothes_desc" gorm:"type:varchar(500);null" binding:"omitempty,max=500"`
		ClothesMerk  *string    `json:"clothes_merk" gorm:"type:varchar(75);null" binding:"omitempty,max=75"`
		ClothesColor string     `json:"clothes_color" gorm:"type:varchar(36);not null" binding:"required,max=36"`
		ClothesPrice *int       `json:"clothes_price" gorm:"type:int;null" binding:"omitempty,min=1,max=99999999"`
		ClothesBuyAt *time.Time `json:"clothes_buy_at" gorm:"type:date;null" binding:"required"`
		ClothesQty   int        `json:"clothes_qty" gorm:"type:int;not null"`
		ClothesImage *string    `json:"clothes_image" gorm:"type:varchar(1000);null"`
		IsFaded      bool       `json:"is_faded" gorm:"type:boolean;not null" binding:"required"`
		HasWashed    bool       `json:"has_washed" gorm:"type:boolean;not null" binding:"required"`
		HasIroned    bool       `json:"has_ironed" gorm:"type:boolean;not null" binding:"required"`
		IsFavorite   bool       `json:"is_favorite" gorm:"type:boolean;not null" binding:"required"`
		IsScheduled  bool       `json:"is_scheduled" gorm:"type:boolean;not null"`
		CreatedAt    time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
		UpdatedAt    *time.Time `json:"updated_at" gorm:"type:timestamp;null"`
		DeletedAt    *time.Time `json:"deleted_at" gorm:"type:timestamp;null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Dictionary
		ClothesMadeFrom    string     `json:"clothes_made_from" gorm:"type:varchar(36);not null" binding:"required,min=36"`
		DictionaryMadeFrom Dictionary `json:"-" gorm:"foreignKey:ClothesMadeFrom;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesType        string     `json:"clothes_type" gorm:"type:varchar(36);not null" binding:"required,min=36"`
		DictionaryType     Dictionary `json:"-" gorm:"foreignKey:ClothesType;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesCategory    string     `json:"clothes_category" gorm:"type:varchar(36);not null" binding:"required,min=36"`
		DictionaryCategory Dictionary `json:"-" gorm:"foreignKey:ClothesCategory;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesSize        string     `json:"clothes_size" gorm:"type:varchar(3);not null" binding:"required,min=3"`
		DictionarySize     Dictionary `json:"-" gorm:"foreignKey:ClothesSize;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ClothesGender      string     `json:"clothes_gender" gorm:"type:varchar(6);not null" binding:"required,min=6"`
		DictionaryGender   Dictionary `json:"-" gorm:"foreignKey:ClothesGender;references:DictionaryName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	ClothesHeader struct {
		ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesName  string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		ClothesColor string    `json:"clothes_color" gorm:"type:varchar(36);not null"`
		ClothesQty   int       `json:"clothes_qty" gorm:"type:int;not null"`
		ClothesImage *string   `json:"clothes_image" gorm:"type:varchar(1000);null"`
		IsFaded      bool      `json:"is_faded" gorm:"type:boolean;not null"`
		HasWashed    bool      `json:"has_washed" gorm:"type:boolean;not null"`
		HasIroned    bool      `json:"has_ironed" gorm:"type:boolean;not null"`
		IsFavorite   bool      `json:"is_favorite" gorm:"type:boolean;not null"`
		IsScheduled  bool      `json:"is_scheduled" gorm:"type:boolean;not null"`
		// FK - Dictionary
		ClothesType     string `json:"clothes_type" gorm:"type:varchar(36);not null"`
		ClothesCategory string `json:"clothes_category" gorm:"type:varchar(36);not null"`
		ClothesSize     string `json:"clothes_size" gorm:"type:varchar(3);not null"`
		ClothesGender   string `json:"clothes_gender" gorm:"type:varchar(6);not null"`
	}
	ClothesDetail struct {
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
		// FK - Dictionary
		ClothesMadeFrom string `json:"clothes_made_from" gorm:"type:varchar(36);not null"`
		ClothesType     string `json:"clothes_type" gorm:"type:varchar(36);not null"`
		ClothesCategory string `json:"clothes_category" gorm:"type:varchar(36);not null"`
		ClothesSize     string `json:"clothes_size" gorm:"type:varchar(3);not null"`
		ClothesGender   string `json:"clothes_gender" gorm:"type:varchar(6);not null"`
	}
	ClothesPlanDestroy struct {
		ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesName     string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
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
	ClothesDeleted struct {
		ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		ClothesName  string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		ClothesImage *string   `json:"clothes_image" gorm:"type:varchar(1000);null"`
		ClothesQty   int       `json:"clothes_qty" gorm:"type:int;not null"`
		ClothesColor string    `json:"clothes_color" gorm:"type:varchar(36);not null"`
		DeletedAt    time.Time `json:"deleted_at" gorm:"type:timestamp;null"`
		// FK - Dictionary
		ClothesType     string `json:"clothes_type" gorm:"type:varchar(36);not null"`
		ClothesCategory string `json:"clothes_category" gorm:"type:varchar(36);not null"`
		ClothesSize     string `json:"clothes_size" gorm:"type:varchar(3);not null"`
		ClothesGender   string `json:"clothes_gender" gorm:"type:varchar(6);not null"`
	}
	SchedulerClothesUnused struct {
		ClothesName     string    `json:"clothes_name" gorm:"type:varchar(36);not null"`
		TotalUsed       int       `json:"total_used" gorm:"type:int;not null"`
		LastUsed        time.Time `json:"last_used" gorm:"type:timestamp;null"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		// FK - Dictionary
		ClothesType string `json:"clothes_type" gorm:"type:varchar(36);not null"`
	}
	SchedulerClothesUnironed struct {
		ClothesName     string  `json:"clothes_name" gorm:"type:varchar(36);not null"`
		HasWashed       bool    `json:"has_washed" gorm:"type:boolean;not null"`
		IsFavorite      bool    `json:"is_favorite" gorm:"type:boolean;not null"`
		IsScheduled     bool    `json:"is_scheduled" gorm:"type:boolean;not null"`
		Username        string  `json:"username" gorm:"type:varchar(36);not null"`
		TelegramUserId  *string `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool    `json:"telegram_is_valid"`
		// FK - Dictionary
		ClothesMadeFrom string `json:"clothes_made_from" gorm:"type:varchar(36);not null"`
		ClothesType     string `json:"clothes_type" gorm:"type:varchar(36);not null"`
	}
)
