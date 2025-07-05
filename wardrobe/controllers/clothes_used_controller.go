package controllers

import (
	"errors"
	"net/http"
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
}

func NewClothesUsedController(clothesUsedService services.ClothesUsedService) *ClothesUsedController {
	return &ClothesUsedController{ClothesUsedService: clothesUsedService}
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
