package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedClothesUseds(repo repositories.ClothesUsedRepository, userRepo repositories.UserRepository, clothesRepo repositories.ClothesRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	for i := 0; i < count; i++ {
		user, err := userRepo.FindOneRandom()
		if err != nil {
			log.Printf("failed to seed clothes used %d: %v\n", i, err)
		}

		clothes, err := clothesRepo.FindOneRandom()
		if err != nil {
			log.Printf("failed to seed clothes used %d: %v\n", i, err)
		}

		clothesUsed := factories.ClothesUsedFactory(clothes.ID)

		err = repo.CreateClothesUsed(&clothesUsed, user.ID)
		if err != nil {
			log.Printf("failed to seed clothes used %d: %v\n", i, err)
		}
		success++
	}
	log.Printf("Seeder : Success to seed %d Clothes Used", success)
}
