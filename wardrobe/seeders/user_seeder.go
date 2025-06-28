package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedUsers(repo repositories.UserRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	for i := 0; i < count; i++ {
		user := factories.UserFactory()
		_, err := repo.CreateUser(&user)
		if err != nil {
			log.Printf("failed to seed user %d: %v\n", i, err)
		}
		success++
	}
	log.Printf("Seeder : Success to seed %d User", success)
}
