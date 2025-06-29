package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedOutfitUseds(repo repositories.OutfitUsedRepository, userRepo repositories.UserRepository,
	outfitRepo repositories.OutfitRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	user, err := userRepo.FindOneHasOutfitAndClothesRandom()
	if err != nil {
		log.Printf("failed to seed outfit used %v\n", err)
	}

	for _, dt := range user {
		for i := 0; i < count; i++ {
			outfit, err := outfitRepo.FindOneRandom(dt.ID)
			if err != nil {
				log.Printf("failed to seed outfit used user - %s at idx - %d : %v\n", dt.Username, i, err)
			}

			outfitUsed := factories.OutfitUsedFactory(outfit.ID)
			err = repo.CreateOutfitUsed(&outfitUsed, dt.ID)
			if err != nil {
				log.Printf("failed to seed outfit used user - %s at idx - %d : %v\n", dt.Username, i, err)
			}
			success++
		}
	}

	log.Printf("Seeder : Success to seed %d Outfit Used", success)
}
