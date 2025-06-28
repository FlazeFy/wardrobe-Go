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
	userTrackRepo repositories.UserTrackRepository,
) {
	seeders.SeedAdmins(adminRepo, 5)
	seeders.SeedUsers(userRepo, 20)
	seeders.SeedDictionaries(dictionaryRepo)
	seeders.SeedQuestions(questionRepo, 15)
	seeders.SeedFeedbacks(feedbackRepo, userRepo, 10)
	seeders.SeedHistories(historyRepo, userRepo, 40)
	seeders.SeedUserTracks(userTrackRepo, userRepo, 15)
}
