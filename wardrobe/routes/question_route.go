package routes

import (
	"wardrobe/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouteQuestion(api *gin.RouterGroup, questionController *controllers.QuestionController) {
	// Public Routes
	question := api.Group("/questions")
	{
		question.POST("/", questionController.CreateQuestion)
		question.GET("/", questionController.GetAllQuestion)
	}
}
