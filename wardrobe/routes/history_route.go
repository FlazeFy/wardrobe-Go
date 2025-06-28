package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteHistory(api *gin.RouterGroup, historyController *controllers.HistoryController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient))
	history := protected.Group("/histories")
	{
		history.GET("/", historyController.GetAllHistory)
		history.DELETE("/destroy/:id", historyController.HardDeleteHistoryById, middleware.AuditTrailMiddleware(db, "hard_delete_history_by_id"))
	}
}
