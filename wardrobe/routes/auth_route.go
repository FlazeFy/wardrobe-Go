package routes

import (
	"wardrobe/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouteAuth(api *gin.RouterGroup, authController *controllers.AuthController) {
	// Public Routes
	auth := api.Group("/auths")
	{
		auth.POST("/register", authController.BasicRegister)
		auth.POST("/login", authController.BasicLogin)
		auth.POST("/signout", authController.BasicSignOut)
	}
}
