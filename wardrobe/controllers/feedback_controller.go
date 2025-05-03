package controllers

import (
	"net/http"
	"wardrobe/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FeedbackController struct {
	DB *gorm.DB
}

func NewFeedbackController(db *gorm.DB) *FeedbackController {
	return &FeedbackController{DB: db}
}

// Queries
func (c *FeedbackController) GetAllFeedback(ctx *gin.Context) {
	// Models
	var data []models.Feedback

	// Query
	c.DB.Find(&data)

	// Response
	status := http.StatusNotFound
	var res interface{} = nil

	if len(data) > 0 {
		status = http.StatusOK
		res = data
	}

	ctx.JSON(status, gin.H{
		"data":    res,
		"message": "feedback fetched",
	})
}

// Command
func (c *FeedbackController) CreateFeedback(ctx *gin.Context) {
	// Models
	var req models.Feedback

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	// Query : Add Feedback
	feedback := models.Feedback{
		FeedbackRate: req.FeedbackRate,
		FeedbackBody: req.FeedbackBody,
	}
	if err := c.DB.Create(&feedback).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create feedback",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"data":    feedback,
		"message": "feedback created",
	})
}

func (c *FeedbackController) HardDeleteFeedbackById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Models
	var feedback models.Feedback

	// Query
	result := c.DB.Unscoped().First(&feedback, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "feedback not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&feedback)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "feedback permanentally deleted",
	})
}
