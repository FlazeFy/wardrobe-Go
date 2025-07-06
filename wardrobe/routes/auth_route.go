package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetUpRouteAuth(api *gin.RouterGroup, authController *controllers.AuthController, redisClient *redis.Client) {
	// Public Routes
	auth := api.Group("/auths")
	{
		// User
		auth.POST("/register", authController.BasicRegister)

		// All Role
		auth.POST("/login", authController.BasicLogin)
		auth.GET("/google", authController.GoogleLogin)
		auth.GET("/google/callback", authController.GoogleRegister)
	}

	// Protected Routes - All Role
	protectedAll := api.Group("/")
	protectedAll.Use(middleware.AuthMiddleware(redisClient, "user", "admin"))
	authAll := protectedAll.Group("/auths")
	{
		authAll.POST("/signout", authController.BasicSignOut)
	}
}
