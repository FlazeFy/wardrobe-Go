package factories

import (
	"wardrobe/models"
)

func DictionaryFactory(dctName, dctType string) models.Dictionary {
	return models.Dictionary{
		DictionaryName: dctName,
		DictionaryType: dctType,
	}
}
