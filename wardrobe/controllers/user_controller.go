package controllers

import (
	"errors"
	"math"
	"net/http"
	"strings"
	"wardrobe/config"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (c *UserController) GetAllUser(ctx *gin.Context) {
	// Param
	order := ctx.Param("order")
	username := ctx.Param("username")

	// Pagination
	pagination := utils.GetPagination(ctx)

	// Validator : Target Column Validator
	order = strings.ToLower(order)
	if !utils.Contains(config.QueryOrder, order) {
		utils.BuildResponseMessage(ctx, "failed", "user", "order is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Get All User
	res, total, err := c.UserService.GetAllUser(pagination, order, username)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "user", "empty", http.StatusNotFound, nil, nil)
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
	utils.BuildResponseMessage(ctx, "success", "user", "get", http.StatusOK, res, metadata)
}
