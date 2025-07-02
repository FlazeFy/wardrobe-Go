package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteDictionary(api *gin.RouterGroup, dictionaryController *controllers.DictionaryController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes - All Role
	protectedAll := api.Group("/")
	protectedAll.Use(middleware.AuthMiddleware(redisClient, "user", "admin"))
	dictionaryAll := protectedAll.Group("/dictionaries")
	{
		dictionaryAll.GET("/", dictionaryController.GetAllDictionary)
		dictionaryAll.GET("/:dictionary_type", dictionaryController.GetDictionaryByType)
	}

	// Protected Routes - Admin
	protectedAdmin := api.Group("/")
	protectedAdmin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	dictionaryAdmin := protectedAdmin.Group("/dictionaries")
	{
		dictionaryAdmin.POST("/", dictionaryController.CreateDictionary, middleware.AuditTrailMiddleware(db, "post_create_dictionary"))
		dictionaryAdmin.DELETE("/destroy/:id", dictionaryController.HardDeleteDictionaryById, middleware.AuditTrailMiddleware(db, "hard_delete_dictionary_by_id"))
	}
}
