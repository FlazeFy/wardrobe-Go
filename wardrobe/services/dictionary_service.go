package services

import (
	"errors"
	"wardrobe/models"
	"wardrobe/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Dictionary Interface
type DictionaryService interface {
	GetAllDictionary() ([]models.Dictionary, error)
	GetDictionaryByType(dictionaryType string) ([]models.Dictionary, error)
	CreateDictionary(dictionary *models.Dictionary) error
	HardDeleteDictionaryByID(ID uuid.UUID) error
}

// Dictionary Struct
type dictionaryService struct {
	dictionaryRepo repositories.DictionaryRepository
}

// Dictionary Constructor
func NewDictionaryService(dictionaryRepo repositories.DictionaryRepository) DictionaryService {
	return &dictionaryService{
		dictionaryRepo: dictionaryRepo,
	}
}

func (r *dictionaryService) GetAllDictionary() ([]models.Dictionary, error) {
	return r.dictionaryRepo.FindAllDictionary()
}

func (r *dictionaryService) GetDictionaryByType(dictionaryType string) ([]models.Dictionary, error) {
	return r.dictionaryRepo.FindDictionaryByType(dictionaryType)
}

func (r *dictionaryService) CreateDictionary(dictionary *models.Dictionary) error {
	// Repo : Find Dictionary By Type
	_, err := r.dictionaryRepo.FindOneDictionaryByName(dictionary.DictionaryName)
	if err == nil {
		return gorm.ErrDuplicatedKey
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.dictionaryRepo.CreateDictionary(dictionary)
}

func (r *dictionaryService) HardDeleteDictionaryByID(ID uuid.UUID) error {
	return r.dictionaryRepo.HardDeleteDictionaryByID(ID)
}
