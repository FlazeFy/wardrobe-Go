package controllers

import (
	"net/http"
	"strings"
	"wardrobe/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DictionaryController struct {
	DB *gorm.DB
}

func NewDictionaryController(db *gorm.DB) *DictionaryController {
	return &DictionaryController{DB: db}
}

// Queries
func (c *DictionaryController) GetAllDictionary(ctx *gin.Context) {
	// Models
	var data []models.Dictionary

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
		"message": "dictionary fetched",
	})
}

func (c *DictionaryController) GetDictionaryByType(ctx *gin.Context) {
	// Params
	dictionaryType := ctx.Param("dictionary_type")

	// Models
	var data []models.Dictionary

	// Query
	if err := c.DB.Where("dictionary_type = ?", dictionaryType).Find(&data).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "dictionary not found",
		})
		return
	}

	// Response
	status := http.StatusNotFound
	var res interface{} = nil

	if len(data) > 0 {
		status = http.StatusOK
		res = data
	}

	ctx.JSON(status, gin.H{
		"data":    res,
		"message": "dictionary fetched",
	})
}

// Command
func (c *DictionaryController) CreateDictionary(ctx *gin.Context) {
	// Models
	var req models.Dictionary

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	// Validate : Dictionary Type Rules
	allowedTypes := []string{"social_type", "interactions_mood"}
	isValidType := false
	for _, t := range allowedTypes {
		if req.DictionaryType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "dictionary_type must be one of: " + strings.Join(allowedTypes, ", "),
		})
		return
	}

	// Query : Add Dictionary
	dictionary := models.Dictionary{
		DictionaryType: req.DictionaryType,
		DictionaryName: req.DictionaryName,
	}
	if err := c.DB.Create(&dictionary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"data":    dictionary,
		"message": "dictionary created",
	})
}
