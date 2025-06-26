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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "something went wrong",
		})
		return
	}

	// Response
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "feedback not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
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
			"status":  "failed",
			"message": "history not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&history)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "history permanentally deleted",
	})
}
