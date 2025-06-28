package routes

import (
	"wardrobe/repositories"
	"wardrobe/seeders"

	"gorm.io/gorm"
)

func SetUpSeeder(db *gorm.DB, adminRepo repositories.AdminRepository) {
	seeders.SeedAdmins(adminRepo, 5)
}
