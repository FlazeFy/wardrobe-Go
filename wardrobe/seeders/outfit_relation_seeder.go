package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedOutfitRelations(repo repositories.OutfitRelationRepository, userRepo repositories.UserRepository,
	clothesRepo repositories.ClothesRepository, outfitRepo repositories.OutfitRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	user, err := userRepo.FindOneHasOutfitAndClothesRandom()
	if err != nil {
		log.Printf("failed to seed outfit relation %v\n", err)
	}

	for _, dt := range user {
		for i := 0; i < count; i++ {
			clothes, err := clothesRepo.FindOneRandom(dt.ID)
			if err != nil {
				log.Printf("failed to seed outfit relation user - %s at idx - %d : %v\n", dt.Username, i, err)
			}

			outfit, err := outfitRepo.FindOneRandom(dt.ID)
			if err != nil {
				log.Printf("failed to seed outfit relation user - %s at idx - %d : %v\n", dt.Username, i, err)
			}

			outfitRel := factories.OutfitRelationFactory(outfit.ID, clothes.ID)
			err = repo.CreateOutfitRelation(&outfitRel, dt.ID)
			if err != nil {
				log.Printf("failed to seed outfit relation user - %s at idx - %d : %v\n", dt.Username, i, err)
			}
			success++
		}
	}

	log.Printf("Seeder : Success to seed %d Outfit Relation", success)
}
