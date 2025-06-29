package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedWashs(repo repositories.WashRepository, userRepo repositories.UserRepository, clothesRepo repositories.ClothesRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	user, err := userRepo.FindOneHasOutfitAndClothesRandom()
	if err != nil {
		log.Printf("failed to seed wash %v\n", err)
	}

	for _, dt := range user {
		for i := 0; i < count; i++ {
			clothes, err := clothesRepo.FindOneRandom(dt.ID)
			if err != nil {
				log.Printf("failed to seed wash user - %s at idx - %d : %v\n", dt.Username, i, err)
			}

			wash := factories.WashFactory(clothes.ID)
			err = repo.CreateWash(&wash, dt.ID)
			if err != nil {
				log.Printf("failed to seed wash user - %s at idx - %d : %v\n", dt.Username, i, err)
			}
			success++
		}
	}
	log.Printf("Seeder : Success to seed %d Wash", success)
}
