package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteFeedback(api *gin.RouterGroup, feedbackController *controllers.FeedbackController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - User
	protectedUser := api.Group("/")
	protectedUser.Use(middleware.AuthMiddleware(redisClient, "user"))
	feedbackUser := protectedUser.Group("/feedbacks")
	{
		feedbackUser.POST("/", feedbackController.CreateFeedback, middleware.AuditTrailMiddleware(db, "post_crete_feedback"))
	}

	// Protected Routes - Admin
	protectedAdmin := api.Group("/")
	protectedAdmin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	feedbackAdmin := protectedAdmin.Group("/feedbacks")
	{
		feedbackAdmin.GET("/", feedbackController.GetAllFeedback)
		feedbackAdmin.DELETE("/destroy/:id", feedbackController.HardDeleteFeedbackById, middleware.AuditTrailMiddleware(db, "hard_delete_feedback_by_id"))
	}
}
