package routes

import (
	"wardrobe/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouteAuth(api *gin.RouterGroup, authController *controllers.AuthController) {
	// Public Routes
	auth := api.Group("/auths")
	{
		// User
		auth.POST("/register", authController.BasicRegister)

		// All Role
		auth.POST("/login", authController.BasicLogin)
		auth.POST("/signout", authController.BasicSignOut)
		auth.GET("/google", authController.GoogleLogin)
		auth.GET("/google/callback", authController.GoogleRegister)
	}
}
