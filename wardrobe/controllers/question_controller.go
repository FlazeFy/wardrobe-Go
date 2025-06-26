package controllers

import (
	"errors"
	"net/http"
	"wardrobe/models"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionController struct {
	QuestionService services.QuestionService
}

func NewQuestionController(questionService services.QuestionService) *QuestionController {
	return &QuestionController{QuestionService: questionService}
}

// Queries
func (c *QuestionController) GetAllQuestion(ctx *gin.Context) {
	questions, err := c.QuestionService.GetAllQuestion()

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "question", "get", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "question", "get", http.StatusOK, questions, nil)
}

// Command
func (c *QuestionController) CreateQuestion(ctx *gin.Context) {
	// Models
	var req models.Question

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BuildResponseMessage(ctx, "failed", "question", "invalid request body", http.StatusBadRequest, nil, nil)
		return
	}

	// Query : Add Question
	question := models.Question{
		Question: req.Question,
	}
	err := c.QuestionService.CreateQuestion(&question)
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "question", "post", http.StatusCreated, nil, nil)
}
