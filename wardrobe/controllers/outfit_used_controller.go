package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutfitUsedController struct {
	OutfitUsedService services.OutfitUsedService
	StatsService      services.StatsService
}

func NewOutfitUsedController(outfitUsedService services.OutfitUsedService, statsService services.StatsService) *OutfitUsedController {
	return &OutfitUsedController{
		OutfitUsedService: outfitUsedService,
		StatsService:      statsService,
	}
}

func (c *OutfitUsedController) GetMonthlyOutfitUsedByClothesIdAndYear(ctx *gin.Context) {
	// Param
	yearStr := ctx.Param("year")
	clothesId := ctx.Param("clothes_id")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "outfit used", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Validator Clothes Id
	if clothesId != "all" {
		_, err := uuid.Parse(clothesId)
		if err != nil {
			utils.BuildResponseMessage(ctx, "failed", "outfit used", "invalid clothes_id", http.StatusBadRequest, nil, nil)
			return
		}
	}

	// Validator Year
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "outfit used", "invalid year", http.StatusBadRequest, nil, nil)
		return
	}

	// Service: Get Most Context
	clothes, err := c.StatsService.GetMonthlyClothesUsedByClothesIdAndYear(year, "outfit_useds", "created_at", "clothes_id", clothesId, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "outfit used", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "outfit used", "get", http.StatusOK, clothes, nil)
}
