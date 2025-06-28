package routes

import (
	"wardrobe/repositories"
	"wardrobe/seeders"

	"gorm.io/gorm"
)

func SetUpSeeder(
	db *gorm.DB, adminRepo repositories.AdminRepository, userRepo repositories.UserRepository,
	dictionaryRepo repositories.DictionaryRepository,
) {
	seeders.SeedAdmins(adminRepo, 5)
	seeders.SeedUsers(userRepo, 20)
	seeders.SeedDictionaries(dictionaryRepo)
}
