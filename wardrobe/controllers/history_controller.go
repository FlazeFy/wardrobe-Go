package controllers

import (
	"net/http"
	"wardrobe/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HistoryController struct {
	DB *gorm.DB
}

func NewHistoryController(db *gorm.DB) *HistoryController {
	return &HistoryController{DB: db}
}

// Queries
func (c *HistoryController) GetAllHistory(ctx *gin.Context) {
	// Models
	var data []models.History

	// Query
	result := c.DB.Preload("User").Find(&data)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": "history not found",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "history fetched",
	})
}

// Command
func (c *HistoryController) HardDeleteHistoryById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Models
	var history models.History

	// Query
	result := c.DB.Unscoped().First(&history, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "history not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&history)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "history permanentally deleted",
	})
}
