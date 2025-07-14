package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteStats(api *gin.RouterGroup, clothesController *controllers.ClothesController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - User
	protectedUser := api.Group("/")
	protectedUser.Use(middleware.AuthMiddleware(redisClient, "user"))
	statsUser := protectedUser.Group("/stats")
	{
		statsUser.GET("/most_context/:target_col", clothesController.GetMostContextClothes)
	}
}
