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
	protectedUser := api.Group("/")
	protectedUser.Use(middleware.AuthMiddleware(redisClient, "user"))
	scheduleUser := protectedUser.Group("/schedules")
	{
		scheduleUser.GET("/by_day/:day", scheduleController.GetScheduleByDay)
		scheduleUser.GET("/by_tomorrow/:day", scheduleController.GetScheduleForTomorrow)
		scheduleUser.POST("/", scheduleController.CreateSchedule, middleware.AuditTrailMiddleware(db, "post_create_schedule"))
		scheduleUser.DELETE("/destroy/:id", scheduleController.HardDeleteScheduleById, middleware.AuditTrailMiddleware(db, "hard_delete_schedule_by_id"))
	}
}
