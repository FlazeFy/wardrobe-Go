package controllers

import (
	"errors"
	"math"
	"net/http"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ErrorController struct {
	ErrorService services.ErrorService
}

func NewErrorController(errorsService services.ErrorService) *ErrorController {
	return &ErrorController{ErrorService: errorsService}
}

// Queries
func (c *ErrorController) GetAllError(ctx *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(ctx)

	// Service : Get All Error
	errorsList, total, err := c.ErrorService.GetAllError(pagination)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "error", "empty", http.StatusNotFound, nil, nil)
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
	utils.BuildResponseMessage(ctx, "success", "error", "get", http.StatusOK, errorsList, metadata)
}

// Command
func (c *ErrorController) HardDeleteErrorById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Parse Param UUID
	errorID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "error", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Hard Delete Error By ID
	err = c.ErrorService.HardDeleteErrorByID(errorID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "error", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "error", "hard delete", http.StatusOK, nil, nil)
}
