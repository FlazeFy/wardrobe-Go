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

func (c *ClothesContext) GetAllClothesHeader(category, order string, userID uuid.UUID) ([]ClothesHeader, error) {
	// Model
	var clothes []ClothesHeader

	// Ordering Prep
	is_desc := true
	if order == "asc" {
		is_desc = false
	}

	// Model
	query := c.DB.Table("clothes").
		Select("id,clothes_name,clothes_image,clothes_size,clothes_gender,clothes_color,clothes_category,clothes_type,clothes_qty,is_faded,has_washed,has_ironed,is_favorite,is_scheduled").
		Where("created_by = ?", userID).
		Where("deleted_at is null")

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
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(clothes) == 0 {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}

func (c *ClothesContext) GetAllClothesDetail(category, order string, userID uuid.UUID) ([]Clothes, error) {
	// Model
	var clothes []Clothes

	// Ordering Prep
	is_desc := true
	if order == "asc" {
		is_desc = false
	}

	// Model
	query := c.DB.Table("clothes").
		Where("created_by = ?", userID).
		Where("deleted_at is null")

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
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(clothes) == 0 {
		return nil, errors.New("clothes used history not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return clothes, nil
}

func (c *ClothesContext) GetClothesShortInfoById(id uuid.UUID) (*ClothesShortInfo, error) {
	// Model
	var clothes ClothesShortInfo

	// Query
	result := c.DB.Table("clothes").
		Select("clothes_name,clothes_type,clothes_category").
		Where("id = ?", id).
		First(&clothes)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &clothes, nil
}

func (c *ClothesContext) GetClothesLastCreated(ctx string, userID uuid.UUID) (*ClothesLastCreated, error) {
	// Model
	var clothes ClothesLastCreated

	// Query
	result := c.DB.Table("clothes").
		Select("clothes_name, "+ctx).
		Where("created_by = ?", userID).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx},
			Desc:   true,
		}).First(&clothes)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &clothes, nil
}

func (c *ClothesContext) GetClothesLastDeleted(ctx string, userID uuid.UUID) (*ClothesLastDeleted, error) {
	// Model
	var clothes ClothesLastDeleted

	// Query
	result := c.DB.Table("clothes").
		Select("clothes_name, "+ctx).
		Where("created_by = ?", userID).
		Where("deleted_at IS NOT NULL").
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx},
			Desc:   true,
		}).First(&clothes)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("clothes not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &clothes, nil
}

func (c *ClothesContext) GetDeletedClothes(userID uuid.UUID) ([]ClothesDeleted, error) {
	// Model
	var clothes []ClothesDeleted

	// Query
	result := c.DB.Table("clothes").
		Select("id, clothes_name, clothes_image, clothes_size, clothes_gender, clothes_color, clothes_category, clothes_type, clothes_qty, deleted_at").
		Where("created_by = ?", userID).
		Where("deleted_at IS NOT NULL").
		Order("deleted_at DESC").
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

func (c *ClothesContext) GetClothesPlanDestroy(days int) ([]ClothesPlanDestroy, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// Model
	var clothes []ClothesPlanDestroy

	// Query
	result := c.DB.Table("clothes").
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

// Command Scheduler
func (c *ClothesContext) SchedulerHardDeleteClothesById(id uuid.UUID) (int64, error) {
	// Model
	var clothes Clothes

	// Query
	result := c.DB.Unscoped().Where("id", id).Delete(&clothes)

	// Response
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (c *ClothesContext) SchedulerDeleteClothesUsedByClothesId(id uuid.UUID) (int64, error) {
	// Model
	var clothes ClothesUsed

	// Query
	result := c.DB.Unscoped().Where("clothes_id", id).Delete(&clothes)

	// Response
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (c *ClothesContext) SchedulerGetUnusedClothes(days int) ([]SchedulerClothesUnused, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	// Model
	var clothes []SchedulerClothesUnused

	// Query
	result := c.DB.Table("clothes").
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

func (c *ClothesContext) SchedulerGetUnironedClothes() ([]SchedulerClothesUnironed, error) {
	ironable_clothes_made_from := []string{"cotton", "linen", "silk", "rayon"}
	ironable_clothes_type := []string{"pants", "shirt", "jacket", "shorts", "skirt", "dress", "blouse", "sweater", "hoodie", "tie", "coat", "vest", "t-shirt", "jeans", "leggings", "cardigan"}

	// Model
	var clothes []SchedulerClothesUnironed

	// Query
	result := c.DB.Table("clothes").
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
