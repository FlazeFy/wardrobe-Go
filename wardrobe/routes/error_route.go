package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteError(api *gin.RouterGroup, errorController *controllers.ErrorController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - Admin
	protectedAdmin := api.Group("/")
	protectedAdmin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	errorsAdmin := protectedAdmin.Group("/errors")
	{
		errorsAdmin.GET("/", errorController.GetAllError)
	}
}
