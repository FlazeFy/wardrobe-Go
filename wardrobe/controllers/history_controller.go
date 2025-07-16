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

type HistoryController struct {
	HistoryService services.HistoryService
}

func NewHistoryController(historyService services.HistoryService) *HistoryController {
	return &HistoryController{HistoryService: historyService}
}

// Queries
// @Summary      Get All History
// @Description  Returns a list of all users histories
// @Tags         History
// @Accept       json
// @Produce      json
// @Success      200  {object}  others.ResponseGetHistory
// @Failure      400  {object}  others.ResponseBadRequestInvalidUserId
// @Failure      404  {object}  others.ResponseNotFound
// @Failure      500  {object}  others.ResponseInternalServerError
// @Router       /api/v1/histories [get]
func (c *HistoryController) GetAllHistory(ctx *gin.Context) {
	var userID *uuid.UUID
	var role string
	var errRole, errUser error

	// Pagination
	pagination := utils.GetPagination(ctx)

	// Get Role
	role, errRole = utils.GetRole(ctx)
	if errRole != nil {
		utils.BuildResponseMessage(ctx, "failed", "history", errRole.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	if role == "user" {
		userID, errUser = utils.GetUserID(ctx)
		if errUser != nil {
			utils.BuildResponseMessage(ctx, "failed", "history", errUser.Error(), http.StatusBadRequest, nil, nil)
			return
		}
	}

	// Service : Get All History
	var history interface{}
	history, total, err := c.HistoryService.GetAllHistory(pagination, userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "history", "get", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	if role == "user" {
		history = utils.StripFields(history, "username")
	}

	// Response
	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	metadata := gin.H{
		"total":       total,
		"page":        pagination.Page,
		"limit":       pagination.Limit,
		"total_pages": totalPages,
	}
	utils.BuildResponseMessage(ctx, "success", "history", "get", http.StatusOK, history, metadata)
}

// Command
// @Summary      Hard Delete History By Id
// @Description  Permanentally delete history by Id
// @Tags         History
// @Success      200  {object}  others.ResponseHardDeleteHistoryById
// @Failure      400  {object}  others.ResponseBadRequestInvalidUserId
// @Failure      404  {object}  others.ResponseNotFound
// @Failure      500  {object}  others.ResponseInternalServerError
// @Router       /api/v1/histories/destroy/{id} [delete]
// @Param        id  path  string  true  "Id of history"
func (c *HistoryController) HardDeleteHistoryById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Parse Param UUID
	historyID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "history", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "history", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Hard Delete History By ID
	err = c.HistoryService.HardDeleteHistoryByID(historyID, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "history", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "history", "hard delete", http.StatusOK, nil, nil)
}
