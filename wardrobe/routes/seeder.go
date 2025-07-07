package routes

import (
	"wardrobe/repositories"
	"wardrobe/seeders"

	"gorm.io/gorm"
)

func SetUpSeeder(
	db *gorm.DB, adminRepo repositories.AdminRepository, userRepo repositories.UserRepository,
	dictionaryRepo repositories.DictionaryRepository, questionRepo repositories.QuestionRepository,
	feedbackRepo repositories.FeedbackRepository, historyRepo repositories.HistoryRepository,
	userTrackRepo repositories.UserTrackRepository, errorRepo repositories.ErrorRepository,
	clothesRepo repositories.ClothesRepository, clothesUsedRepo repositories.ClothesUsedRepository,
	userWeatherRepo repositories.UserWeatherRepository, outfitRepo repositories.OutfitRepository,
	outfitRelationRepo repositories.OutfitRelationRepository, scheduleRepo repositories.ScheduleRepository,
	outfitUsedRepo repositories.OutfitUsedRepository, washRepo repositories.WashRepository,
) {
	seeders.SeedAdmins(adminRepo, 5)
	seeders.SeedUsers(userRepo, 20)
	seeders.SeedDictionaries(dictionaryRepo)
	seeders.SeedQuestions(questionRepo, 15)
	seeders.SeedFeedbacks(feedbackRepo, userRepo, 10)
	seeders.SeedHistories(historyRepo, userRepo, 40)
	seeders.SeedUserTracks(userTrackRepo, userRepo, 15)
	seeders.SeedErrors(errorRepo, 25)
	seeders.SeedClothes(clothesRepo, userRepo, 200)
	seeders.SeedUserWeathers(userWeatherRepo, userRepo, 100)
	seeders.SeedOutfits(outfitRepo, userRepo, 60)

	// This count per User
	seeders.SeedClothesUseds(clothesUsedRepo, userRepo, clothesRepo, 20)
	seeders.SeedOutfitRelations(outfitRelationRepo, userRepo, clothesRepo, outfitRepo, 10)
	seeders.SeedSchedules(scheduleRepo, userRepo, clothesRepo, 7)
	seeders.SeedOutfitUseds(outfitUsedRepo, userRepo, outfitRepo, 7)
	seeders.SeedWashs(washRepo, userRepo, clothesRepo, 10)
}
