package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteClothes(api *gin.RouterGroup, clothesController *controllers.ClothesController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - User
	protectedUser := api.Group("/")
	protectedUser.Use(middleware.AuthMiddleware(redisClient, "user"))
	clothesUser := protectedUser.Group("/clothes")
	{
		clothesUser.POST("/", clothesController.CreateClothes, middleware.AuditTrailMiddleware(db, "post_create_history"))
		clothesUser.GET("/last_history", clothesController.GetClothesLastHistory)
		clothesUser.DELETE("/destroy/:id", clothesController.HardDeleteClothesById, middleware.AuditTrailMiddleware(db, "hard_delete_clothes_by_id"))
	}

	// Protected Routes - All Role
	protectedAll := api.Group("/")
	protectedAll.Use(middleware.AuthMiddleware(redisClient, "user", "admin"))
	clothesAll := protectedAll.Group("/clothes")
	{
		clothesAll.GET("/header/:category/:order", clothesController.GetAllClothesHeader)
		clothesAll.GET("/detail/:category/:order", clothesController.GetAllClothesDetail)
		clothesAll.GET("/trash", clothesController.GetDeletedClothes)
		clothesAll.PUT("/recover/:id", clothesController.RecoverDeletedClothesById, middleware.AuditTrailMiddleware(db, "put_recover_deleted_clothes_by_id"))
		clothesAll.DELETE("/:id", clothesController.SoftDeleteClothesById, middleware.AuditTrailMiddleware(db, "soft_delete_clothes_by_id"))
	}
}
