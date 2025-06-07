package controllers

import (
	"net/http"
	"wardrobe/models"
	"wardrobe/models/others"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	var req models.User

	// Validate : Request Body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Validate : If username or email already been used
	var existing models.User
	if err := ac.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username or email has already been used",
		})
		return
	}

	// Hashing
	if err := utils.HashPassword(&req, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Query
	user := models.User{
		ID:              uuid.New(),
		Username:        req.Username,
		Password:        req.Password,
		Email:           req.Email,
		TelegramUserId:  req.TelegramUserId,
		TelegramIsValid: false,
	}
	result := ac.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error,
		})
		return
	}

	// JWT Token Generate
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error,
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
			"message": err.Error(),
		})
		return
	}

	// Query
	if err := ac.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid email",
		})
		return
	}

	// Validate : Password
	if err := utils.CheckPassword(&user, loginReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid password",
		})
		return
	}

	// JWT Token Generate
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "error generating token",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User login successfully",
		"token":   token,
	})
}


func (ac *AuthController) SignOut(c *gin.Context) {
	// Header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing authorization header",
			"status":  "failed",
		})
		return
	}

	// Clean Bearer
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid authorization header",
			"status":  "failed",
		})
		return
	}

	// Reset Token By Adding Blacklist Redis
	err := ac.AuthService.SignOut(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "user signout successfully",
		"status":  "success",
	})
}
