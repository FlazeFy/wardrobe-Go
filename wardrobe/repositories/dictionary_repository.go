package repositories

import (
	"strings"
	"time"
	"wardrobe/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Dictionary Interface
type DictionaryRepository interface {
	CreateDictionary(dictionary *models.Dictionary) error
	FindAllDictionary() ([]models.Dictionary, error)
	FindDictionaryByType(dictionaryType string) ([]models.Dictionary, error)
	FindOneDictionaryByName(dictionaryName string) (*models.Dictionary, error)
	HardDeleteDictionaryByID(ID uuid.UUID) error

	// For Seeder
	DeleteAll() error
}

// Dictionary Struct
type dictionaryRepository struct {
	db *gorm.DB
}

// Dictionary Constructor
func NewDictionaryRepository(db *gorm.DB) DictionaryRepository {
	return &dictionaryRepository{db: db}
}

func (r *dictionaryRepository) FindAllDictionary() ([]models.Dictionary, error) {
	// Model
	var dictionaries []models.Dictionary

	// Query
	if err := r.db.Order("dictionary_type ASC").
		Order("dictionary_name ASC").
		Find(&dictionaries).Error; err != nil {
		return nil, err
	}
	if len(dictionaries) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return dictionaries, nil
}

func (r *dictionaryRepository) FindDictionaryByType(dictionaryType string) ([]models.Dictionary, error) {
	// Model
	var dictionaries []models.Dictionary

	// Query
	if err := r.db.Where("dictionary_type", dictionaryType).
		Order("dictionary_type ASC").
		Order("dictionary_name ASC").
		Find(&dictionaries).Error; err != nil {
		return nil, err
	}
	if len(dictionaries) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return dictionaries, nil
}

func (r *dictionaryRepository) FindOneDictionaryByName(dictionaryName string) (*models.Dictionary, error) {
	// Model
	var dictionaries models.Dictionary

	// Prepare Search
	dictionaryNameLower := strings.ToLower(dictionaryName)

	// Query
	if err := r.db.Where("LOWER(dictionary_name) = ?", dictionaryNameLower).
		Order("dictionary_type ASC").
		Order("dictionary_name ASC").
		First(&dictionaries).Error; err != nil {
		return nil, err
	}

	return &dictionaries, nil
}

func (r *dictionaryRepository) CreateDictionary(dictionary *models.Dictionary) error {
	// Default
	dictionary.ID = uuid.New()
	dictionary.CreatedAt = time.Now()

	// Query
	return r.db.Create(dictionary).Error
}

func (r *dictionaryRepository) HardDeleteDictionaryByID(ID uuid.UUID) error {
	// Query
	result := r.db.Unscoped().Where("id = ?", ID).Delete(&models.Dictionary{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// For Seeder
func (r *dictionaryRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.Dictionary{}).Error
}
