package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedSchedules(repo repositories.ScheduleRepository, userRepo repositories.UserRepository, clothesRepo repositories.ClothesRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	user, err := userRepo.FindOneHasOutfitAndClothesRandom()
	if err != nil {
		log.Printf("failed to seed schedule %v\n", err)
	}

	for _, dt := range user {
		for i := 0; i < count; i++ {
			clothes, err := clothesRepo.FindOneRandom(dt.ID)
			if err != nil {
				log.Printf("failed to seed schedule user - %s at idx - %d : %v\n", dt.Username, i, err)
			}

			schedule := factories.ScheduleFactory(clothes.ID)
			err = repo.CreateSchedule(&schedule, dt.ID)
			if err != nil {
				log.Printf("failed to seed schedule user - %s at idx - %d : %v\n", dt.Username, i, err)
			}
			success++
		}
	}
	log.Printf("Seeder : Success to seed %d Schedule", success)
}
