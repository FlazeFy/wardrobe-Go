package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
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
			auth.POST("/signout", authController.SignOut)
		}
		question := api.Group("/question")
		{
			question.POST("/", questionController.CreateQuestion)
			question.GET("/", questionController.GetAllQuestion)
		}

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(redisClient))

		feedback := protected.Group("/feedback")
		{
			feedback.GET("/", feedbackController.GetAllFeedback)
			feedback.POST("/", feedbackController.CreateFeedback, middleware.AuditTrailMiddleware(db, "post_crete_feedback"))
			feedback.DELETE("/destroy/:id", feedbackController.HardDeleteFeedbackById, middleware.AuditTrailMiddleware(db, "hard_delete_feedback_by_id"))
		}
		dictionary := protected.Group("/dictionary")
		{
			dictionary.GET("/", dictionaryController.GetAllDictionary)
			dictionary.GET("/:dictionary_type", dictionaryController.GetDictionaryByType)
			dictionary.POST("/", dictionaryController.CreateDictionary, middleware.AuditTrailMiddleware(db, "post_create_dictionary"))
		}
		history := protected.Group("/history")
		{
			history.GET("/", historyController.GetAllHistory)
			history.DELETE("/destroy/:id", historyController.HardDeleteHistoryById, middleware.AuditTrailMiddleware(db, "hard_delete_history_by_id"))
		}
		clothes := protected.Group("/clothes")
		{
			clothes.POST("/", clothesController.CreateClothes, middleware.AuditTrailMiddleware(db, "post_create_history"))
			clothes.GET("/header/:category/:order", clothesController.GetAllClothesHeader)
			clothes.GET("/detail/:category/:order", clothesController.GetAllClothesDetail)
			clothes.GET("/last_history", clothesController.GetClothesLastHistory)
			clothes.GET("/trash", clothesController.GetDeletedClothes)
			clothes.GET("/history/:clothes_id/:order", clothesController.GetClothesUsedHistory)
			clothes.PUT("/recover/:id", clothesController.RecoverDeletedClothesById, middleware.AuditTrailMiddleware(db, "put_recover_deleted_clothes_by_id"))
			clothes.DELETE("/:id", clothesController.SoftDeleteClothesById, middleware.AuditTrailMiddleware(db, "soft_delete_clothes_by_id"))
			clothes.DELETE("/destroy/:id", clothesController.HardDeleteClothesById, middleware.AuditTrailMiddleware(db, "hard_delete_clothes_by_id"))
			clothes.POST("/history", clothesController.CreateClothesUsed, middleware.AuditTrailMiddleware(db, "post_create_clothes_used"))
			clothes.DELETE("/destroy_used/:id", clothesController.HardDeleteClothesUsedById, middleware.AuditTrailMiddleware(db, "hard_delete_clothes_used_by_id"))
		}
		schedule := protected.Group("/schedule")
		{
			schedule.GET("/by_day/:day", scheduleController.GetScheduleByDay, middleware.AuditTrailMiddleware(db, "get_schedule_by_day"))
			schedule.GET("/by_tomorrow/:day", scheduleController.GetScheduleForTomorrow, middleware.AuditTrailMiddleware(db, "get_schedule_for_tommorow"))
			schedule.POST("/", scheduleController.CreateSchedule, middleware.AuditTrailMiddleware(db, "post_create_schedule"))
			schedule.DELETE("/destroy/:id", scheduleController.HardDeleteScheduleById, middleware.AuditTrailMiddleware(db, "hard_delete_schedule_by_id"))
		}
	}
}
