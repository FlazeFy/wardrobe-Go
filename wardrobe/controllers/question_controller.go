package controllers

import (
	"errors"
	"math"
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
	// Pagination
	pagination := utils.GetPagination(ctx)

	// Service : Get All Question
	questions, total, err := c.QuestionService.GetAllQuestion(pagination)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "question", "empty", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	metadata := gin.H{
		"total":       total,
		"page":        pagination.Page,
		"limit":       pagination.Limit,
		"total_pages": totalPages,
	}
	utils.BuildResponseMessage(ctx, "success", "question", "get", http.StatusOK, questions, metadata)
}

// Command
func (c *QuestionController) CreateQuestion(ctx *gin.Context) {
	// Models
	var req models.Question

	// Validate JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.BuildValidationError(err)
		utils.BuildResponseMessage(ctx, "failed", "question", formattedErrors, http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Create Question
	question := models.Question{
		Question: req.Question,
	}
	err := c.QuestionService.CreateQuestion(&question)
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "question", "post", http.StatusCreated, nil, nil)
}
