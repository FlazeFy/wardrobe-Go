package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteClothesUsed(api *gin.RouterGroup, clothesUsedController *controllers.ClothesUsedController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient))
	clothesUsed := protected.Group("/clothes_used")
	{
		clothesUsed.GET("/history/:clothes_id/:order", clothesUsedController.GetClothesUsedHistory)
		clothesUsed.POST("/history", clothesUsedController.CreateClothesUsed, middleware.AuditTrailMiddleware(db, "post_create_clothes_used"))
		clothesUsed.DELETE("/destroy_used/:id", clothesUsedController.HardDeleteClothesUsedById, middleware.AuditTrailMiddleware(db, "hard_delete_clothes_used_by_id"))
	}
}
