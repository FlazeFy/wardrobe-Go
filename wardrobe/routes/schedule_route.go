package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteSchedule(api *gin.RouterGroup, scheduleController *controllers.ScheduleController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient))
	schedule := protected.Group("/schedules")
	{
		schedule.GET("/by_day/:day", scheduleController.GetScheduleByDay)
		schedule.GET("/by_tomorrow/:day", scheduleController.GetScheduleForTomorrow)
		schedule.POST("/", scheduleController.CreateSchedule, middleware.AuditTrailMiddleware(db, "post_create_schedule"))
		schedule.DELETE("/destroy/:id", scheduleController.HardDeleteScheduleById, middleware.AuditTrailMiddleware(db, "hard_delete_schedule_by_id"))
	}
}
