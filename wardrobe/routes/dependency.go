package routes

import (
	"sync"
	"wardrobe/bots/line"
	"wardrobe/bots/telegram"
	"wardrobe/cache"
	"wardrobe/controllers"
	"wardrobe/repositories"
	"wardrobe/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpDependency(r *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	// Migrate DB
	MigrateAll(db)

	// Dependency Cache
	statsCache := cache.NewStatsCache(redisClient)

	// Dependency Repositories
	adminRepo := repositories.NewAdminRepository(db)
	clothesRepo := repositories.NewClothesRepository(db)
	clothesUsedRepo := repositories.NewClothesUsedRepository(db)
	dictionaryRepo := repositories.NewDictionaryRepository(db)
	errorRepo := repositories.NewErrorRepository(db)
	feedbackRepo := repositories.NewFeedbackRepository(db)
	historyRepo := repositories.NewHistoryRepository(db)
	outfitRelationRepo := repositories.NewOutfitRelationRepository(db)
	outfitRepo := repositories.NewOutfitRepository(db)
	questionRepo := repositories.NewQuestionRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	userRepo := repositories.NewUserRepository(db)
	userWeatherRepo := repositories.NewUserWeatherRepository(db)
	washRepo := repositories.NewWashRepository(db)
	userTrackRepo := repositories.NewUserTrackRepository(db)
	outfitUsedRepo := repositories.NewOutfitUsedRepository(db)
	googleTokenRepo := repositories.NewGoogleTokenRepository(db)
	statsRepo := repositories.NewStatsRepository(db)

	// Dependency Services
	adminService := services.NewAdminService(adminRepo)
	authService := services.NewAuthService(userRepo, adminRepo, googleTokenRepo, redisClient)
	clothesService := services.NewClothesService(clothesRepo, userRepo, scheduleRepo, clothesUsedRepo, washRepo, outfitRelationRepo)
	clothesUsedService := services.NewClothesUsedService(clothesUsedRepo)
	dictionaryService := services.NewDictionaryService(dictionaryRepo)
	errorService := services.NewErrorService(errorRepo)
	feedbackService := services.NewFeedbackService(feedbackRepo)
	historyService := services.NewHistoryService(historyRepo)
	questionService := services.NewQuestionService(questionRepo)
	scheduleService := services.NewScheduleService(scheduleRepo, userRepo, clothesRepo)
	userService := services.NewUserService(userRepo)
	userWeatherService := services.NewUserWeatherService(userWeatherRepo)
	statsService := services.NewStatsService(statsRepo, redisClient, statsCache)
	washService := services.NewWashService(washRepo)

	// Dependency Controller
	authController := controllers.NewAuthController(authService)
	clothesController := controllers.NewClothesController(clothesService, statsService)
	clothesUsedController := controllers.NewClothesUsedController(clothesUsedService)
	dictionaryController := controllers.NewDictionaryController(dictionaryService)
	feedbackController := controllers.NewFeedbackController(feedbackService)
	historyController := controllers.NewHistoryController(historyService)
	questionController := controllers.NewQuestionController(questionService)
	scheduleController := controllers.NewScheduleController(scheduleService)
	errorController := controllers.NewErrorController(errorService)

	// Routes Endpoint
	SetUpRoutes(r, db, redisClient, authController, questionController, feedbackController,
		dictionaryController, historyController, clothesController, clothesUsedController, scheduleController,
		errorController,
	)

	// Task Scheduler
	SetUpScheduler(adminService, errorService, historyService, clothesService, clothesUsedService,
		scheduleService, washService, questionService, userService, userWeatherService,
	)

	// Seeder & Factories
	SetUpSeeder(db, adminRepo, userRepo, dictionaryRepo, questionRepo, feedbackRepo, historyRepo,
		userTrackRepo, errorRepo, clothesRepo, clothesUsedRepo, userWeatherRepo, outfitRepo,
		outfitRelationRepo, scheduleRepo, outfitUsedRepo, washRepo,
	)

	// Line Bot
	r.POST("/api/v1/callback/line", line.LineHandler())

	// Telegram Bot
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		telegram.TelegramHandler()
	}()
}
