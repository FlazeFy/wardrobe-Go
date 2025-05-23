package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	questionController := controllers.NewQuestionController(db)
	feedbackController := controllers.NewFeedbackController(db)
	dictionaryController := controllers.NewDictionaryController(db)
	historyController := controllers.NewHistoryController(db)
	clothesController := controllers.NewClothesController(db)
	scheduleController := controllers.NewScheduleController(db)

	api := r.Group("/api/v2")
	{
		// Public Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
		question := api.Group("/question")
		{
			question.POST("/", questionController.CreateQuestion)
			question.GET("/", questionController.GetAllQuestion)
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
		clothes := protected.Group("/clothes")
		{
			clothes.POST("/", clothesController.CreateClothes)
			clothes.GET("/header/:category/:order", clothesController.GetAllClothesHeader)
			clothes.GET("/detail/:category/:order", clothesController.GetAllClothesDetail)
			clothes.GET("/last_history", clothesController.GetClothesLastHistory)
			clothes.GET("/trash", clothesController.GetDeletedClothes)
			clothes.GET("/history/:clothes_id/:order", clothesController.GetClothesUsedHistory)
			clothes.PUT("/recover/:id", clothesController.RecoverDeletedClothesById)
			clothes.DELETE("/:id", clothesController.SoftDeleteClothesById)
			clothes.DELETE("/destroy/:id", clothesController.HardDeleteClothesById)
			clothes.POST("/history", clothesController.CreateClothesUsed)
			clothes.DELETE("/destroy_used/:id", clothesController.HardDeleteClothesUsedById)
		}
		schedule := protected.Group("/schedule")
		{
			schedule.GET("/by_day/:day", scheduleController.GetScheduleByDay)
			schedule.GET("/by_tomorrow/:day", scheduleController.GetScheduleForTomorrow)
			schedule.POST("/", scheduleController.CreateSchedule)
			schedule.DELETE("/destroy/:id", scheduleController.HardDeleteScheduleById)
		}
	}
}
