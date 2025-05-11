package controllers

import (
	"net/http"
	"strings"
	"wardrobe/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	result := c.DB.Where("dictionary_type = ?", dictionaryType).Find(&data)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": "something went wrong",
		})
		return
	}

	// Response
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": "no dictionary found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    data,
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
	allowedTypes := []string{"used_context", "wash_type", "clothes_type", "clothes_category", "clothes_size", "clothes_made_from", "clothes_gender"}
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

	// Query : Check Dictionary Name
	var existing models.Dictionary
	if err := c.DB.Where("dictionary_name = ?", req.DictionaryName).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "dictionary with the same name already exists",
		})
		return
	}

	// Query : Add Dictionary
	dictionary := models.Dictionary{
		ID:             uuid.New(),
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
