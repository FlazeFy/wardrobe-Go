package seeders

import (
	"log"
	"wardrobe/factories"
	"wardrobe/repositories"
)

func SeedQuestions(repo repositories.QuestionRepository, count int) {
	// Empty Table
	repo.DeleteAll()

	// Fill Table
	var success = 0
	for i := 0; i < count; i++ {
		question := factories.QuestionFactory()
		err := repo.CreateQuestion(&question)
		if err != nil {
			log.Printf("failed to seed question %d: %v\n", i, err)
		}
		success++
	}
	log.Printf("Seeder : Success to seed %d Question", success)
}
