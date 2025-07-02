package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteClothesUsed(api *gin.RouterGroup, clothesUsedController *controllers.ClothesUsedController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - User
	protectedUser := api.Group("/")
	protectedUser.Use(middleware.AuthMiddleware(redisClient, "user"))
	clothesUsedUser := protectedUser.Group("/clothes_used")
	{
		clothesUsedUser.GET("/history/:clothes_id/:order", clothesUsedController.GetClothesUsedHistory)
		clothesUsedUser.POST("/history", clothesUsedController.CreateClothesUsed, middleware.AuditTrailMiddleware(db, "post_create_clothes_used"))
		clothesUsedUser.DELETE("/destroy_used/:id", clothesUsedController.HardDeleteClothesUsedById, middleware.AuditTrailMiddleware(db, "hard_delete_clothes_used_by_id"))
	}
}
