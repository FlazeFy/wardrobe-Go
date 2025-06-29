package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedOutfits(repo repositories.OutfitRepository, userRepo repositories.UserRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	for i := 0; i < count; i++ {
		outfit := factories.OutfitFactory()
		user, err := userRepo.FindOneRandom()
		if err != nil {
			log.Printf("failed to seed outfit %d: %v\n", i, err)
		}

		err = repo.CreateOutfit(&outfit, user.ID)
		if err != nil {
			log.Printf("failed to seed outfit %d: %v\n", i, err)
		}
		success++
	}
	log.Printf("Seeder : Success to seed %d Outfit", success)
}
