package main

import (
	"log"
	"os"
	"wardrobe/config"
	"wardrobe/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load Env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading ENV")
	}

	// Connect DB
	db := config.ConnectDatabase()

	// Setup Gin
	router := gin.Default()

	// Setup Redis
	redisClient := config.InitRedis()

	// Dependency For Routes, Migrations, Seeders, and Task Scheduler
	routes.SetUpDependency(router, db, redisClient)

	// Run server
	port := os.Getenv("PORT")
	router.Run(":" + port)

	log.Printf("Wardrobe is running on port %s\n", port)
}
