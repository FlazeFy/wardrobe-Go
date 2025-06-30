package controllers

import (
	"errors"
	"net/http"
	"wardrobe/models"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedbackController struct {
	FeedbackService services.FeedbackService
}

func NewFeedbackController(feedbackService services.FeedbackService) *FeedbackController {
	return &FeedbackController{FeedbackService: feedbackService}
}

// Queries
func (c *FeedbackController) GetAllFeedback(ctx *gin.Context) {
	// Service : Get All Feedback
	feedback, err := c.FeedbackService.GetAllFeedback()

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "feedback", "get", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "feedback", "get", http.StatusOK, feedback, nil)
}

// Command
func (c *FeedbackController) CreateFeedback(ctx *gin.Context) {
	// Models
	var req models.Feedback

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BuildResponseMessage(ctx, "failed", "feedback", "invalid request body", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "feedback", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Add Feedback
	err = c.FeedbackService.CreateFeedback(&req, *userID)
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "feedback", "post", http.StatusCreated, nil, nil)
}

func (c *FeedbackController) HardDeleteFeedbackById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Parse Param UUID
	feedbackID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "feedback", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Hard Delete Feedback By ID
	err = c.FeedbackService.HardDeleteFeedbackByID(feedbackID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "feedback", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "feedback", "hard delete", http.StatusOK, nil, nil)
}
