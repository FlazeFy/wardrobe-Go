package services

import (
	"errors"
	"wardrobe/models"
	"wardrobe/repositories"

	"gorm.io/gorm"
)

// Dictionary Interface
type DictionaryService interface {
	GetAllDictionary() ([]models.Dictionary, error)
	GetDictionaryByType(dictionaryType string) ([]models.Dictionary, error)
	CreateDictionary(dictionary *models.Dictionary) error
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
		return errors.New("dictionary already exist")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.dictionaryRepo.CreateDictionary(dictionary)
}
