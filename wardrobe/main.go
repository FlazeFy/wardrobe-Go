package main

import (
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	redisClient := config.InitRedis()

	routes.SetUpDependency(router, db, redisClient)

	// Run server
	router.Run(":9000")
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
		&models.AuditTrail{},
	)
}
