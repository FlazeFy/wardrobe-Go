package controllers

import (
	"errors"
	"net/http"
	"wardrobe/config"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WashController struct {
	WashService  services.WashService
	StatsService services.StatsService
}

func NewWashController(washService services.WashService, statsService services.StatsService) *WashController {
	return &WashController{
		WashService:  washService,
		StatsService: statsService,
	}
}

// Query
func (c *WashController) GetMostContextWash(ctx *gin.Context) {
	// Param
	targetCol := ctx.Param("target_col")

	// Validator : Target Column Validator
	if !utils.Contains(config.StatsWashField, targetCol) {
		utils.BuildResponseMessage(ctx, "failed", "wash", "target_col is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "wash", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service: Get Most Context
	wash, err := c.StatsService.GetMostUsedContext("washes", targetCol, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "wash", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "wash", "get", http.StatusOK, wash, nil)
}
