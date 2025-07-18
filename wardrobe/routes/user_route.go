package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteUser(api *gin.RouterGroup, userController *controllers.UserController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - Admin
	protectedAdmin := api.Group("/")
	protectedAdmin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	usersAdmin := protectedAdmin.Group("/users")
	{
		usersAdmin.GET("/:order/:username", userController.GetAllUser)
	}
}
