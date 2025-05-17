package controllers

import (
	"net/http"
	"wardrobe/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionController struct {
	DB *gorm.DB
}

func NewQuestionController(db *gorm.DB) *QuestionController {
	return &QuestionController{DB: db}
}

// Queries
func (c *QuestionController) GetAllQuestion(ctx *gin.Context) {
	// Models
	var data []models.Question

	// Query
	result := c.DB.Find(&data)

	// Response
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "question not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "failed",
		"data":    data,
		"message": "question fetched",
	})
}

// Command
func (c *QuestionController) CreateQuestion(ctx *gin.Context) {
	// Models
	var req models.Question

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request body",
		})
		return
	}

	// Query : Add Question
	question := models.Question{
		ID:       uuid.New(),
		Question: req.Question,
		Answer:   nil,
		IsShow:   false,
	}
	if err := c.DB.Create(&question).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to create question",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    question,
		"message": "question created",
	})
}
