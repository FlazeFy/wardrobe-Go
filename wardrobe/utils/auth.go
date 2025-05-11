package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(ctx *gin.Context) (*uuid.UUID, error) {
	userIDValue, exists := ctx.Get("userId")
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
		return nil, errors.New("invalid user id type")
	}
}
