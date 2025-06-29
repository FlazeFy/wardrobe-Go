package routes

import (
	"wardrobe/controllers"
	"wardrobe/repositories"
	"wardrobe/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpDependency(r *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	// Dependency Repositories
	adminRepo := repositories.NewAdminRepository(db)
	clothesRepo := repositories.NewClothesRepository(db)
	clothesUsedRepo := repositories.NewClothesUsedRepository(db)
	dictionaryRepo := repositories.NewDictionaryRepository(db)
	errorRepo := repositories.NewErrorRepository(db)
	feedbackRepo := repositories.NewFeedbackRepository(db)
	historyRepo := repositories.NewHistoryRepository(db)
	// outfitRelationRepo := repositories.NewOutfitRelationRepository(db)
	questionRepo := repositories.NewQuestionRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	userRepo := repositories.NewUserRepository(db)
	userWeatherRepo := repositories.NewUserWeatherRepository(db)
	washRepo := repositories.NewWashRepository(db)
	userTrackRepo := repositories.NewUserTrackRepository(db)

	// Dependency Services
	adminService := services.NewAdminService(adminRepo)
	authService := services.NewAuthService(userRepo, adminRepo, redisClient)
	clothesService := services.NewClothesService(clothesRepo, userRepo)
	clothesUsedService := services.NewClothesUsedService(clothesUsedRepo)
	dictionaryService := services.NewDictionaryService(dictionaryRepo)
	errorService := services.NewErrorService(errorRepo)
	feedbackService := services.NewFeedbackService(feedbackRepo)
	historyService := services.NewHistoryService(historyRepo)
	questionService := services.NewQuestionService(questionRepo)
	scheduleService := services.NewScheduleService(scheduleRepo, userRepo, clothesRepo)
	userService := services.NewUserService(userRepo)
	userWeatherService := services.NewUserWeatherService(userWeatherRepo)
	washService := services.NewWashService(washRepo)

	// Dependency Controller
	authController := controllers.NewAuthController(authService)
	clothesController := controllers.NewClothesController(clothesService)
	clothesUsedController := controllers.NewClothesUsedController(clothesUsedService)
	dictionaryController := controllers.NewDictionaryController(dictionaryService)
	feedbackController := controllers.NewFeedbackController(feedbackService)
	historyController := controllers.NewHistoryController(historyService)
	questionController := controllers.NewQuestionController(questionService)
	scheduleController := controllers.NewScheduleController(scheduleService)

	// Routes Endpoint
	SetUpRoutes(r, db, redisClient,
		authController, questionController, feedbackController,
		dictionaryController, historyController, clothesController,
		clothesUsedController, scheduleController,
	)

	// Task Scheduler
	SetUpScheduler(
		adminService, errorService, historyService,
		clothesService, clothesUsedService, scheduleService,
		washService, questionService, userService,
		userWeatherService,
	)

	// Seeder & Factories
	SetUpSeeder(db, adminRepo, userRepo, dictionaryRepo, questionRepo, feedbackRepo, historyRepo,
		userTrackRepo, errorRepo, clothesRepo, clothesUsedRepo)
}
