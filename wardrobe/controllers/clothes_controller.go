package controllers

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClothesController struct {
	ClothesService services.ClothesService
}

func NewClothesController(clothesService services.ClothesService) *ClothesController {
	return &ClothesController{ClothesService: clothesService}
}

// Query
func (c *ClothesController) GetClothesLastHistory(ctx *gin.Context) {
	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Get Clothes Last History
	data, err := c.ClothesService.GetClothesLastHistory(*userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "clothes", "empty", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes", "get", http.StatusOK, data, nil)
}

func (c *ClothesController) GetDeletedClothes(ctx *gin.Context) {
	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Get Deleted Clothes
	res, err := c.ClothesService.GetDeletedClothes(*userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "clothes", "empty", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes", "get", http.StatusOK, res, nil)
}

func (c *ClothesController) GetAllClothesHeader(ctx *gin.Context) {
	// Param
	category := ctx.Param("category")
	order := ctx.Param("order")

	// Pagination
	pagination := utils.GetPagination(ctx)

	if category != "all" {
		// Validator Contain : Clothes Category
		if !utils.Contains(config.ClothesCategories, category) {
			utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes category is not valid", http.StatusBadRequest, nil, nil)
			return
		}
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Get All Clothes Header
	res, total, err := c.ClothesService.GetAllClothesHeader(pagination, category, order, *userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "clothes", "empty", http.StatusNotFound, nil, nil)
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
	utils.BuildResponseMessage(ctx, "success", "clothes", "get", http.StatusOK, res, metadata)
}

func (c *ClothesController) GetAllClothesDetail(ctx *gin.Context) {
	// Param
	category := ctx.Param("category")
	order := ctx.Param("order")

	// Pagination
	pagination := utils.GetPagination(ctx)

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Get All Clothes Detail
	res, total, err := c.ClothesService.GetAllClothesDetail(pagination, category, order, *userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "clothes", "empty", http.StatusNotFound, nil, nil)
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
	utils.BuildResponseMessage(ctx, "success", "clothes", "get", http.StatusOK, res, metadata)
}

// Command
func (c *ClothesController) CreateClothes(ctx *gin.Context) {
	// Mandatory Field
	clothesName := ctx.PostForm("clothes_name")
	clothesColor := ctx.PostForm("clothes_color")
	clothesMadeFrom := ctx.PostForm("clothes_made_from")
	clothesType := ctx.PostForm("clothes_type")
	clothesCategory := ctx.PostForm("clothes_category")
	clothesSize := ctx.PostForm("clothes_size")
	clothesGender := ctx.PostForm("clothes_gender")
	clothesQty, _ := strconv.Atoi(ctx.PostForm("clothes_qty"))
	isFaded, _ := strconv.ParseBool(ctx.PostForm("is_faded"))
	hasWashed, _ := strconv.ParseBool(ctx.PostForm("has_washed"))
	hasIroned, _ := strconv.ParseBool(ctx.PostForm("has_ironed"))
	isFavorite, _ := strconv.ParseBool(ctx.PostForm("is_favorite"))
	isScheduled, _ := strconv.ParseBool(ctx.PostForm("is_scheduled"))

	// Optional Field
	var clothesDesc *string
	desc := ctx.PostForm("clothes_desc")
	if desc != "" {
		clothesDesc = &desc
	}
	var clothesMerk *string
	merk := ctx.PostForm("clothes_merk")
	if merk != "" {
		clothesMerk = &merk
	}

	var clothesPrice *int
	priceStr := ctx.PostForm("clothes_price")
	if priceStr != "" {
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes price is not valid", http.StatusBadRequest, nil, nil)
			return
		}
		clothesPrice = &price
	}

	var clothesBuyAt *time.Time
	buyAtStr := ctx.PostForm("clothes_buy_at")
	if buyAtStr != "" {
		datetime, err := time.Parse(time.RFC3339, buyAtStr)
		if err != nil {
			utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes buy at is not valid", http.StatusBadRequest, nil, nil)
			return
		}
		clothesBuyAt = &datetime
	}

	// Validator Contain : Clothes Category
	if !utils.Contains(config.ClothesCategories, clothesCategory) {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes category is not valid", http.StatusBadRequest, nil, nil)
		return
	}
	// Validator Contain : Clothes Gender
	if !utils.Contains(config.ClothesGenders, clothesGender) {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes gender is not valid", http.StatusBadRequest, nil, nil)
		return
	}
	// Validator Contain : Clothes Type
	if !utils.Contains(config.ClothesTypes, clothesType) {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes type is not valid", http.StatusBadRequest, nil, nil)
		return
	}
	// Validator Contain : Clothes Made From
	if !utils.Contains(config.ClothesMadeFroms, clothesMadeFrom) {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes made from is not valid", http.StatusBadRequest, nil, nil)
		return
	}
	// Validator Contain : Clothes Made From
	if !utils.Contains(config.ClothesMadeFroms, clothesMadeFrom) {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes made from is not valid", http.StatusBadRequest, nil, nil)
		return
	}
	// Validator Contain : Clothes Size
	if !utils.Contains(config.ClothesSizes, clothesSize) {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "clothes size is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	clothes := &models.Clothes{
		ClothesName:     clothesName,
		ClothesDesc:     clothesDesc,
		ClothesMerk:     clothesMerk,
		ClothesColor:    clothesColor,
		ClothesPrice:    clothesPrice,
		ClothesBuyAt:    clothesBuyAt,
		ClothesQty:      clothesQty,
		ClothesImage:    nil,
		IsFaded:         isFaded,
		HasWashed:       hasWashed,
		HasIroned:       hasIroned,
		IsFavorite:      isFavorite,
		IsScheduled:     isScheduled,
		ClothesMadeFrom: clothesMadeFrom,
		ClothesType:     clothesType,
		ClothesCategory: clothesCategory,
		ClothesSize:     clothesSize,
		ClothesGender:   clothesGender,
	}

	// Service : Create Clothes
	clothes, err = c.ClothesService.CreateClothes(clothes, *userID)
	if err != nil {
		if err.Error() == "clothes already exist" {
			utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusConflict, nil, nil)
			return
		}
		if err.Error() == "user contact not found" || err.Error() == "clothes not found" {
			utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusNotFound, nil, nil)
			return
		}

		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes", "post", http.StatusCreated, clothes, nil)
}

func (c *ClothesController) SoftDeleteClothesById(ctx *gin.Context) {
	// Param
	id := ctx.Param("id")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Parse Param UUID
	clothesID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Soft Delete Clothes By Id
	err = c.ClothesService.SoftDeleteClothesById(*userID, clothesID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "clothes", "empty", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes", "soft delete", http.StatusOK, nil, nil)
}

func (c *ClothesController) RecoverDeletedClothesById(ctx *gin.Context) {
	// Param
	id := ctx.Param("id")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Parse Param UUID
	clothesID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Recover Deleted Clothes By Id
	err = c.ClothesService.RecoverDeletedClothesById(*userID, clothesID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "clothes", "empty", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes", "recover", http.StatusOK, nil, nil)
}

func (c *ClothesController) HardDeleteClothesById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Parse Param UUID
	clothesID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Hard Delete Clothes By ID
	err = c.ClothesService.HardDeleteClothesById(clothesID, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "clothes", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "clothes", "hard delete", http.StatusOK, nil, nil)
}
