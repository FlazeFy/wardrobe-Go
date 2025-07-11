package controllers

import (
	"errors"
	"net/http"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/models/others"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
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
		formattedErrors := utils.BuildValidationError(err)
		utils.BuildResponseMessage(ctx, "failed", "register", formattedErrors, http.StatusBadRequest, nil, nil)
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

	utils.BuildResponseMessage(ctx, "success", "user", "register", http.StatusCreated, gin.H{
		"token": token,
	}, nil)
}

func (c *AuthController) GoogleLogin(ctx *gin.Context) {
	url := config.GetGoogleOAuthConfig().AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	ctx.Redirect(http.StatusFound, url)
}

func (c *AuthController) GoogleRegister(ctx *gin.Context) {
	// Validator
	code := ctx.Query("code")
	if code == "" {
		utils.BuildResponseMessage(ctx, "failed", "register", "cant find code", http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Google Register
	token, err := c.AuthService.GoogleRegister(code)
	if err != nil {
		if err.Error() == "username or email has already been used" {
			utils.BuildResponseMessage(ctx, "failed", "register", err.Error(), http.StatusConflict, nil, nil)
			return
		}

		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "user", "register", http.StatusCreated, gin.H{
		"token": token,
	}, nil)
}

func (c *AuthController) BasicLogin(ctx *gin.Context) {
	// Models
	var req others.LoginRequest

	// Validate JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.BuildValidationError(err)
		utils.BuildResponseMessage(ctx, "failed", "auth", formattedErrors, http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Basic Login
	token, role, err := c.AuthService.BasicLogin(req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BuildResponseMessage(ctx, "failed", "auth", "account not found", http.StatusNotFound, nil, nil)
			return
		}
		if err.Error() == "invalid password" {
			utils.BuildResponseMessage(ctx, "failed", "auth", err.Error(), http.StatusBadRequest, nil, nil)
			return
		}

		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", *role, "login", http.StatusOK, gin.H{
		"token": token,
		"role":  role,
	}, nil)
}

func (c *AuthController) BasicSignOut(ctx *gin.Context) {
	// Header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		utils.BuildResponseMessage(ctx, "failed", "auth", "missing authorization header", http.StatusBadRequest, nil, nil)
		return
	}

	// Get Role
	role, errRole := utils.GetRole(ctx)
	if errRole != nil {
		utils.BuildResponseMessage(ctx, "failed", "auth", errRole.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Basic Sign Out
	err := c.AuthService.BasicSignOut(authHeader)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "auth", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	utils.BuildResponseMessage(ctx, "success", role, "sign out", http.StatusOK, nil, nil)
}

func (c *AuthController) GetMyProfile(ctx *gin.Context) {
	// Get Role
	role, err := utils.GetRole(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "profile", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "profile", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Get My Profile
	profile, err := c.AuthService.GetMyProfile(*userID, role)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.BuildResponseMessage(ctx, "failed", "profile", "get", http.StatusNotFound, nil, nil)
		default:
			utils.BuildErrorMessage(ctx, err.Error())
		}
		return
	}

	utils.BuildResponseMessage(ctx, "success", "profile", "get", http.StatusOK, profile, nil)
}
