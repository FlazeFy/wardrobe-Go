package routes

import (
	"wardrobe/controllers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, db *gorm.DB, redisClient *redis.Client,
	authController *controllers.AuthController,
	questionController *controllers.QuestionController,
	feedbackController *controllers.FeedbackController,
	dictionaryController *controllers.DictionaryController,
	historyController *controllers.HistoryController,
	clothesController *controllers.ClothesController,
	clothesUsedController *controllers.ClothesUsedController,
	scheduleController *controllers.ScheduleController,
	errorController *controllers.ErrorController,
	washController *controllers.WashController,
	userWeatherController *controllers.UserWeatherController) {

	// V1 Endpoint
	api := r.Group("/api/v1")

	// Routes Endpoint
	SetUpRouteAuth(api, authController, redisClient)
	SetUpRouteQuestion(api, questionController)
	SetUpRouteFeedback(api, feedbackController, redisClient, db)
	SetUpRouteDictionary(api, dictionaryController, redisClient, db)
	SetUpRouteHistory(api, historyController, redisClient, db)
	SetUpRouteSchedule(api, scheduleController, redisClient, db)
	SetUpRouteClothes(api, clothesController, redisClient, db)
	SetUpRouteClothesUsed(api, clothesUsedController, redisClient, db)
	SetUpRouteError(api, errorController, redisClient, db)
	SetUpRouteStats(api, clothesController, clothesUsedController, scheduleController, washController, userWeatherController, redisClient, db)
}
