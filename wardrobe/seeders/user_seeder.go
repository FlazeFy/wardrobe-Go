package seeders

import (
	"log"
	"os"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedUsers(repo repositories.UserRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	var success = 0

	userUsername := os.Getenv("USER_USERNAME")
	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")
	userTelegramID := os.Getenv("USER_TELEGRAM_USER_ID")
	userTest := factories.UserFactory(&userUsername, &userEmail, &userTelegramID, &userPassword, true)
	_, err := repo.CreateUser(&userTest)
	if err != nil {
		log.Printf("failed to seed admin %d:\n", err)
	}
	success++

	// Fill Table
	for i := 0; i < count; i++ {
		user := factories.UserFactory(nil, nil, nil, nil, false)
		_, err := repo.CreateUser(&user)
		if err != nil {
			log.Printf("failed to seed user %d: %v\n", i, err)
		}
		success++
	}
	log.Printf("Seeder : Success to seed %d User", success)
}
