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
	outfitUsedController *controllers.OutfitUsedController, redisClient *redis.Client, db *gorm.DB) {
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
		statsUserMonthly := statsUser.Group("/monthly")
		{
			statsUserMonthly.GET("/clothes_used/:clothes_id/:year", clothesUsedController.GetMonthlyClothesUsedByClothesIdAndYear)
			statsUserMonthly.GET("/wash/:clothes_id/:year", washController.GetMonthlyWashByClothesIdAndYear)
			statsUserMonthly.GET("/outfit_used/:clothes_id/:year", outfitUsedController.GetMonthlyOutfitUsedByClothesIdAndYear)
		}
	}

	// Protected Routes - Admin
	protectedAdmin := api.Group("/")
	protectedAdmin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	statsAdmin := protectedAdmin.Group("/stats")
	{
		statsAdminMostContext := statsAdmin.Group("/most_context")
		{
			statsAdminMostContext.GET("/clothes/:target_col/:user_id", clothesController.GetMostContextClothesByAdmin)
			statsAdminMostContext.GET("/clothes_used/:target_col/:user_id", clothesUsedController.GetMostContextClothesUsedsByAdmin)
			statsAdminMostContext.GET("/schedule/:target_col/:user_id", scheduleController.GetMostContextScheduleByAdmin)
			statsAdminMostContext.GET("/wash/:target_col/:user_id", washController.GetMostContextWashByAdmin)
		}
	}
}
