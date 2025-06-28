package seeders

import (
	"log"
	"wardrobe/config"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedDictionaries(repo repositories.DictionaryRepository) {
	// Empty Table
	repo.DeleteAll()

	var seedData = []struct {
		DictionaryType  string
		DictionaryNames []string
	}{
		{"clothes_type", config.ClothesTypes},
		{"clothes_category", config.ClothesCategories},
		{"used_context", config.UsedContexts},
		{"clothes_gender", config.ClothesGenders},
		{"clothes_made_from", config.ClothesMadeFroms},
		{"clothes_size", config.ClothesSizes},
		{"day", config.Days},
	}

	// Fill Table
	var success = 0
	for _, dt := range seedData {
		for _, dictionaryName := range dt.DictionaryNames {
			dct := factories.DictionaryFactory(dictionaryName, dt.DictionaryType)
			err := repo.CreateDictionary(&dct)
			if err != nil {
				log.Printf("failed to seed dictionary %s/%s: %v\n", dt.DictionaryType, dictionaryName, err)
			}
			success++
		}
	}
	log.Printf("Seeder : Success to seed %d Dictionary", success)
}
