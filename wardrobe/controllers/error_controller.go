package controllers

import (
	"errors"
	"math"
	"net/http"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
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
			utils.BuildResponseMessage(ctx, "failed", "error", "get", http.StatusNotFound, nil, nil)
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
