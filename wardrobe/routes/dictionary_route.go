package routes

import (
	"wardrobe/controllers"
	middleware "wardrobe/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteDictionary(api *gin.RouterGroup, dictionaryController *controllers.DictionaryController, redisClient *redis.Client, db *gorm.DB) {
	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient))
	dictionary := protected.Group("/dictionaries")
	{
		dictionary.GET("/", dictionaryController.GetAllDictionary)
		dictionary.GET("/:dictionary_type", dictionaryController.GetDictionaryByType)
		dictionary.POST("/", dictionaryController.CreateDictionary, middleware.AuditTrailMiddleware(db, "post_create_dictionary"))
	}
}
