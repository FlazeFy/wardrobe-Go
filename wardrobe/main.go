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

	db := config.ConnectDatabase()
	MigrateAll(db)

	router := gin.Default()
	MigrateAll(db)
	routes.SetUpRoutes(router, db)
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
	)
}
