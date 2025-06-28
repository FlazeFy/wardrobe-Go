package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteClothes(api *gin.RouterGroup, clothesController *controllers.ClothesController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient))
	clothes := protected.Group("/clothes")
	{
		clothes.POST("/", clothesController.CreateClothes, middleware.AuditTrailMiddleware(db, "post_create_history"))
		clothes.GET("/header/:category/:order", clothesController.GetAllClothesHeader)
		clothes.GET("/detail/:category/:order", clothesController.GetAllClothesDetail)
		clothes.GET("/last_history", clothesController.GetClothesLastHistory)
		clothes.GET("/trash", clothesController.GetDeletedClothes)
		clothes.PUT("/recover/:id", clothesController.RecoverDeletedClothesById, middleware.AuditTrailMiddleware(db, "put_recover_deleted_clothes_by_id"))
		clothes.DELETE("/:id", clothesController.SoftDeleteClothesById, middleware.AuditTrailMiddleware(db, "soft_delete_clothes_by_id"))
		clothes.DELETE("/destroy/:id", clothesController.HardDeleteClothesById, middleware.AuditTrailMiddleware(db, "hard_delete_clothes_by_id"))
	}
}
