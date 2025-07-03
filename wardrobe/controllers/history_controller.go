package controllers

import (
	"errors"
	"net/http"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HistoryController struct {
	HistoryService services.HistoryService
}

func NewHistoryController(historyService services.HistoryService) *HistoryController {
	return &HistoryController{HistoryService: historyService}
}

// Queries
func (c *HistoryController) GetAllHistory(ctx *gin.Context) {
	// Service : Get All History
	history, err := c.HistoryService.GetAllHistory()
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "history", "get", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "history", "get", http.StatusOK, history, nil)
}

// Command
func (c *HistoryController) HardDeleteHistoryById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Parse Param UUID
	historyID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "history", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "history", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Hard Delete History By ID
	err = c.HistoryService.HardDeleteHistoryByID(historyID, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "history", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "history", "hard delete", http.StatusOK, nil, nil)
}
