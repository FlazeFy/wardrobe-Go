package main

import (
	"log"
	"os"
	"time"
	"wardrobe/config"
	"wardrobe/routes"

	_ "wardrobe/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Wardrobe API
// @version 1.0
// @description API Documentation for Wardrobe BE - Go Gin
// @host localhost:9000
// @BasePath /api/v1
// @contact.name   Leonardho R. Sitanggang
// @contact.email  flazen.edu@gmail.com

func initLogging() {
	now := time.Now()
	logFileName := "logs/wardrobe-" + now.Format("January-2006") + ".log"

	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	initLogging()
	log.Println("Wardrobe API service is starting...")

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

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	port := os.Getenv("PORT")
	router.Run(":" + port)

	log.Printf("Wardrobe is running on port %s\n", port)
}
