package utils

import (
	"errors"
	"wardrobe/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserContext struct {
	DB *gorm.DB
}

func NewUserContext(db *gorm.DB) *UserContext {
	return &UserContext{DB: db}
}

func GetUserID(ctx *gin.Context) (*uuid.UUID, error) {
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	switch v := userIDValue.(type) {
	case string:
		userID, err := uuid.Parse(v)
		if err != nil {
			return nil, err
		}
		return &userID, nil
	case uuid.UUID:
		return &v, nil
	default:
		return nil, errors.New("invalid user id")
	}
}

func GetRole(c *gin.Context) (string, error) {
	roleVal, exists := c.Get("role")
	if !exists {
		return "", errors.New("role not found in context")
	}

	role, ok := roleVal.(string)
	if !ok {
		return "", errors.New("invalid role format in context")
	}

	return role, nil
}

func (c *UserContext) GetAdminContact() ([]models.UserContact, error) {
	// Model
	var contact []models.UserContact

	// Query
	result := c.DB.Table("admins").
		Select("username, email, telegram_user_id, telegram_is_valid").
		Find(&contact)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(contact) == 0 {
		return nil, errors.New("admin contact not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return contact, nil
}
