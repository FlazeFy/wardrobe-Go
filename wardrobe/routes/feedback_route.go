package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteFeedback(api *gin.RouterGroup, feedbackController *controllers.FeedbackController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient))
	feedback := protected.Group("/feedbacks")
	{
		feedback.GET("/", feedbackController.GetAllFeedback)
		feedback.POST("/", feedbackController.CreateFeedback, middleware.AuditTrailMiddleware(db, "post_crete_feedback"))
		feedback.DELETE("/destroy/:id", feedbackController.HardDeleteFeedbackById, middleware.AuditTrailMiddleware(db, "hard_delete_feedback_by_id"))
	}
}
