package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	feedbackController := controllers.NewFeedbackController(db)
	dictionaryController := controllers.NewDictionaryController(db)
	historyController := controllers.NewHistoryController(db)

	api := r.Group("/api/v2")
	{
		// Public Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())

		feedback := protected.Group("/feedback")
		{
			feedback.GET("/", feedbackController.GetAllFeedback)
			feedback.POST("/", feedbackController.CreateFeedback)
			feedback.DELETE("/destroy/:id", feedbackController.HardDeleteFeedbackById)
		}
		dictionary := protected.Group("/dictionary")
		{
			dictionary.GET("/", dictionaryController.GetAllDictionary)
			dictionary.GET("/:dictionary_type", dictionaryController.GetDictionaryByType)
			dictionary.POST("/", dictionaryController.CreateDictionary)
		}
		history := protected.Group("/history")
		{
			history.GET("/", historyController.GetAllHistory)
			history.DELETE("/destroy/:id", historyController.HardDeleteHistoryById)
		}
	}
}
