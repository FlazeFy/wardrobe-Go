package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedClothes(repo repositories.ClothesRepository, userRepo repositories.UserRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	for i := 0; i < count; i++ {
		clothes := factories.ClothesFactory()
		user, err := userRepo.FindOneRandom()
		if err != nil {
			log.Printf("failed to seed clothes %d: %v\n", i, err)
		}

		_, err = repo.CreateClothes(&clothes, user.ID)
		if err != nil {
			log.Printf("failed to seed clothes %d: %v\n", i, err)
		}
		success++
	}
	log.Printf("Seeder : Success to seed %d Clothes", success)
}
