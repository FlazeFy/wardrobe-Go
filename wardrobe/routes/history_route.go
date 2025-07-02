package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteHistory(api *gin.RouterGroup, historyController *controllers.HistoryController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - User
	protectedUser := api.Group("/")
	protectedUser.Use(middleware.AuthMiddleware(redisClient, "user"))
	historyUser := protectedUser.Group("/histories")
	{
		historyUser.DELETE("/destroy/:id", historyController.HardDeleteHistoryById, middleware.AuditTrailMiddleware(db, "hard_delete_history_by_id"))
	}

	// Protected Routes - All Role
	protectedAll := api.Group("/")
	protectedAll.Use(middleware.AuthMiddleware(redisClient, "admin", "user"))
	historyAll := protectedAll.Group("/histories")
	{
		historyAll.GET("/", historyController.GetAllHistory)
	}
}
