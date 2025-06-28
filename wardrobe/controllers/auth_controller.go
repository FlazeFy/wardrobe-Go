package controllers

import (
	"net/http"
	"wardrobe/models"
	"wardrobe/models/others"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// Command
func (c *AuthController) BasicRegister(ctx *gin.Context) {
	// Models
	var req models.User

	// Validate JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BuildResponseMessage(ctx, "failed", "register", "invalid request body", http.StatusBadRequest, nil, nil)
		return
	}

	// Validate Field
	if req.Username == "" {
		utils.BuildResponseMessage(ctx, "failed", "register", "username is required", http.StatusBadRequest, nil, nil)
		return
	}
	if req.Password == "" {
		utils.BuildResponseMessage(ctx, "failed", "register", "password is required", http.StatusBadRequest, nil, nil)
		return
	}
	if req.Email == "" {
		utils.BuildResponseMessage(ctx, "failed", "register", "email is required", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Basic Register
	token, err := c.AuthService.BasicRegister(req)
	if err != nil {
		if err.Error() == "username or email has already been used" {
			utils.BuildResponseMessage(ctx, "failed", "register", err.Error(), http.StatusConflict, nil, nil)
			return
		}

		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "user", "register", http.StatusCreated, gin.H{
		"token": token,
	}, nil)
}

func (c *AuthController) BasicLogin(ctx *gin.Context) {
	// Models
	var req others.LoginRequest

	// Validate JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BuildResponseMessage(ctx, "failed", "auth", "invalid request body", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Basic Login
	token, err := c.AuthService.BasicLogin(req)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "auth", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "user", "login", http.StatusOK, gin.H{
		"token": token,
	}, nil)
}

func (c *AuthController) BasicSignOut(ctx *gin.Context) {
	// Header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		utils.BuildResponseMessage(ctx, "failed", "auth", "missing authorization header", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Basic Sign Out
	err := c.AuthService.BasicSignOut(authHeader)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "auth", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "user", "sign out", http.StatusOK, nil, nil)
}
