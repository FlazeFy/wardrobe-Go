package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClothesUsedController struct {
	ClothesUsedService services.ClothesUsedService
	StatsService       services.StatsService
}

func NewClothesUsedController(clothesUsedService services.ClothesUsedService, statsService services.StatsService) *ClothesUsedController {
	return &ClothesUsedController{
		ClothesUsedService: clothesUsedService,
		StatsService:       statsService,
	}
}

func (c *ClothesUsedController) GetClothesUsedHistory(ctx *gin.Context) {
	// Params
	clothes_id_param := ctx.Param("clothes_id")
	order := ctx.Param("order")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Clothes Id
	var clothes_id uuid.UUID
	if clothes_id_param == "all" {
		clothes_id = uuid.Nil
	} else {
		clothes_id, err = uuid.Parse(clothes_id_param)
		if err != nil {
			utils.BuildResponseMessage(ctx, "failed", "clothes used", "invalid id", http.StatusBadRequest, nil, nil)
			return
		}
	}

	// Service : Get Clothes Used History
	res, err := c.ClothesUsedService.GetClothesUsedHistory(*userID, clothes_id, order)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes used", "get", http.StatusOK, res, nil)
}

func (c *ClothesUsedController) HardDeleteClothesUsedById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Parse Param UUID
	clothesUsedID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Hard Delete Clothes Used By ID
	err = c.ClothesUsedService.HardDeleteClothesUsedByID(clothesUsedID, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes used", "hard delete", http.StatusOK, nil, nil)
}

func (c *ClothesUsedController) CreateClothesUsed(ctx *gin.Context) {
	// Models
	var req models.ClothesUsed

	// Validate JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.BuildValidationError(err)
		utils.BuildResponseMessage(ctx, "failed", "clothes used", formattedErrors, http.StatusBadRequest, nil, nil)
		return
	}

	// Validator Contain : Used Context
	if !utils.Contains(config.UsedContexts, req.UsedContext) {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "used context is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Create Clothes Used
	clothes_used := models.ClothesUsed{
		ClothesNote: req.ClothesNote,
		ClothesId:   req.ClothesId,
		UsedContext: req.UsedContext,
	}
	err = c.ClothesUsedService.CreateClothesUsed(&req, *userID)
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes used", "post", http.StatusCreated, clothes_used, nil)
}

func (c *ClothesUsedController) GetMostContextClothesUseds(ctx *gin.Context) {
	// Param
	targetCol := ctx.Param("target_col")

	// Validator : Target Column Validator
	if !utils.Contains(config.StatsClothesUsedsField, targetCol) {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "target_col is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service: Get Most Context
	clothes, err := c.StatsService.GetMostUsedContext("clothes_useds", targetCol, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "clothes used", "get", http.StatusOK, clothes, nil)
}

func (c *ClothesUsedController) GetMostContextClothesUsedsByAdmin(ctx *gin.Context) {
	// Param
	targetCol := ctx.Param("target_col")
	userIDStr := ctx.Param("user_id")

	// Validator : Target Column Validator
	if !utils.Contains(config.StatsClothesUsedsField, targetCol) {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "target_col is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Parse Param UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "invalid user id", http.StatusBadRequest, nil, nil)
		return
	}

	// Service: Get Most Context
	clothes, err := c.StatsService.GetMostUsedContext("clothes_useds", targetCol, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "clothes used", "get", http.StatusOK, clothes, nil)
}

func (c *ClothesUsedController) GetMonthlyClothesUsedByClothesIdAndYear(ctx *gin.Context) {
	// Param
	yearStr := ctx.Param("year")
	clothesId := ctx.Param("clothes_id")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Validator Clothes Id
	if clothesId != "all" {
		_, err := uuid.Parse(clothesId)
		if err != nil {
			utils.BuildResponseMessage(ctx, "failed", "clothes used", "invalid clothes_id", http.StatusBadRequest, nil, nil)
			return
		}
	}

	// Validator Year
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "invalid year", http.StatusBadRequest, nil, nil)
		return
	}

	// Service: Get Most Context
	clothes, err := c.StatsService.GetMonthlyClothesUsedByClothesIdAndYear(year, "clothes_useds", "created_at", "clothes_id", clothesId, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "clothes used", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "clothes used", "get", http.StatusOK, clothes, nil)
}
