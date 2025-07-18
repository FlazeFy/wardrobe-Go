package seeders

import (
	"log"
	"os"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedAdmins(repo repositories.AdminRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	var success = 0

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminTelegramID := os.Getenv("ADMIN_TELEGRAM_USER_ID")
	adminTest := factories.AdminFactory(&adminUsername, &adminEmail, &adminTelegramID, &adminPassword, true)
	err := repo.Create(&adminTest)
	if err != nil {
		log.Printf("failed to seed admin %d:\n", err)
	}
	success++

	for i := 0; i < count; i++ {
		admin := factories.AdminFactory(nil, nil, nil, nil, false)
		err := repo.Create(&admin)
		if err != nil {
			log.Printf("failed to seed admin %d: %v\n", i, err)
		}
		success++
	}
	log.Printf("Seeder : Success to seed %d Admin", success)
}
