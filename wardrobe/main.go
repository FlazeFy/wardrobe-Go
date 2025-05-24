package main

import (
	"time"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/routes"
	"wardrobe/schedulers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func main() {
	// Load Env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading ENV")
	}

	// Connect DB
	db := config.ConnectDatabase()
	MigrateAll(db)

	// Setup Gin
	router := gin.Default()
	routes.SetUpRoutes(router, db)

	// Task Scheduler
	c := cron.New()
	Scheduler(c)
	c.Start()
	defer c.Stop()

	// Run server
	router.Run(":9000")
}

func Scheduler(c *cron.Cron) {
	// For Production
	// Clean Scheduler
	c.AddFunc("10 0 * * *", schedulers.SchedulerWeatherRoutineFetch)
	c.AddFunc("0 2 * * *", schedulers.SchedulerCleanHistory)
	c.AddFunc("0 2 * * *", schedulers.SchedulerCleanDeletedClothes)
	c.AddFunc("0 1 * * 1", schedulers.SchedulerAuditError)
	c.AddFunc("20 2 * * 1,3,6", schedulers.SchedulerReminderUnansweredQuestion)
	c.AddFunc("20 2 * * 0,2,5", schedulers.SchedulerReminderUnusedClothes)
	c.AddFunc("20 3 * * *", schedulers.SchedulerReminderUnironedClothes)

	// For Development
	go func() {
		time.Sleep(5 * time.Second)

		// Clean Scheduler
		// schedulers.SchedulerCleanHistory()
		// schedulers.SchedulerCleanDeletedClothes()

		// Audit Scheduler
		// schedulers.SchedulerAuditError()

		// Reminder Scheduler
		// schedulers.SchedulerReminderUnansweredQuestion()
		// schedulers.SchedulerReminderUnusedClothes()
		// schedulers.SchedulerReminderUnironedClothes()

		// Weather Scheduler
		schedulers.SchedulerWeatherRoutineFetch()
	}()
}

func MigrateAll(db *gorm.DB) {
	db.AutoMigrate(
		&models.Admin{},
		&models.User{},
		&models.Dictionary{},
		&models.Error{},
		&models.History{},
		&models.Feedback{},
		&models.Clothes{},
		&models.ClothesUsed{},
		&models.Outfit{},
		&models.OutfitRelation{},
		&models.OutfitUsed{},
		&models.Question{},
		&models.Schedule{},
		&models.UserRequest{},
		&models.UserTrack{},
		&models.UserWeather{},
		&models.Wash{},
	)
}
