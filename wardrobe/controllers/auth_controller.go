package controllers

import (
	"net/http"
	"wardrobe/models"
	"wardrobe/models/others"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// Command
func (ac *AuthController) Register(c *gin.Context) {
	// Models
	var user models.User

	// Validate : Request Body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Validate : If username or email already been used
	var existing models.User
	if err := ac.DB.Where("username = ? OR email = ?", user.Username, user.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username or email has already been used",
		})
		return
	}

	// Hashing
	if err := utils.HashPassword(&user, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Query
	result := ac.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error,
		})
		return
	}

	// JWT Token Generate
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error,
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	// Models
	var loginReq others.LoginRequest
	var user models.User

	// Validate : Request Body
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Query
	if err := ac.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "invalid email",
		})
		return
	}

	// Validate : Password
	if err := utils.CheckPassword(&user, loginReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "invalid password",
		})
		return
	}

	// JWT Token Generate
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "error generating token",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User login successfully",
		"token":   token,
	})
}
