package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteStats(api *gin.RouterGroup, clothesController *controllers.ClothesController, clothesUsedController *controllers.ClothesUsedController,
	scheduleController *controllers.ScheduleController, washController *controllers.WashController, userWeatherController *controllers.UserWeatherController,
	redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - User
	protectedUser := api.Group("/")
	protectedUser.Use(middleware.AuthMiddleware(redisClient, "user"))
	statsUser := protectedUser.Group("/stats")
	{
		statsUserMostContext := statsUser.Group("/most_context")
		{
			statsUserMostContext.GET("/clothes/:target_col", clothesController.GetMostContextClothes)
			statsUserMostContext.GET("/clothes_used/:target_col", clothesUsedController.GetMostContextClothesUseds)
			statsUserMostContext.GET("/schedule/:target_col", scheduleController.GetMostContextSchedule)
			statsUserMostContext.GET("/wash/:target_col", washController.GetMostContextWash)
			statsUserMostContext.GET("/user_weather/:target_col", userWeatherController.GetMostContextUserWeather)
		}
	}
}
