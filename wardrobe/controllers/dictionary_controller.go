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

type DictionaryController struct {
	DictionaryService services.DictionaryService
}

func NewDictionaryController(dictionaryService services.DictionaryService) *DictionaryController {
	return &DictionaryController{DictionaryService: dictionaryService}
}

// Queries
func (c *DictionaryController) GetAllDictionary(ctx *gin.Context) {
	// Service : Get All Dictionary
	dictionary, err := c.DictionaryService.GetAllDictionary()
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "dictionary", "get", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "dictionary", "get", http.StatusOK, dictionary, nil)
}

func (c *DictionaryController) GetDictionaryByType(ctx *gin.Context) {
	// Params
	dictionaryType := ctx.Param("dictionary_type")

	// Validator Contain : Dictionary Type
	if !utils.Contains(config.DictionaryTypes, dictionaryType) {
		utils.BuildResponseMessage(ctx, "failed", "dictionary", "dictionary_type is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Get Dictionary By Type
	dictionary, err := c.DictionaryService.GetDictionaryByType(dictionaryType)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "dictionary", "get", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "dictionary", "get", http.StatusOK, dictionary, nil)
}

// Command
func (c *DictionaryController) CreateDictionary(ctx *gin.Context) {
	// Models
	var req models.Dictionary

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BuildResponseMessage(ctx, "failed", "question", "invalid request body", http.StatusBadRequest, nil, nil)
		return
	}

	// Validator Contain : Dictionary Type
	if !utils.Contains(config.DictionaryTypes, req.DictionaryType) {
		utils.BuildResponseMessage(ctx, "failed", "dictionary", "dictionary_type is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Create Dictionary
	dictionary := models.Dictionary{
		DictionaryType: req.DictionaryType,
		DictionaryName: req.DictionaryName,
	}
	err := c.DictionaryService.CreateDictionary(&dictionary)
	if err != nil {
		if err.Error() == "dictionary already exists" {
			utils.BuildResponseMessage(ctx, "failed", "dictionary", "already exists", http.StatusConflict, nil, nil)
			return
		}

		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "dictionary", "post", http.StatusCreated, nil, nil)
}

func (c *DictionaryController) HardDeleteDictionaryById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Parse Param UUID
	dictionaryID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "dictionary", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Hard Delete Dictionary By ID
	err = c.DictionaryService.HardDeleteDictionaryByID(dictionaryID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "dictionary", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "dictionary", "hard delete", http.StatusOK, nil, nil)
}
