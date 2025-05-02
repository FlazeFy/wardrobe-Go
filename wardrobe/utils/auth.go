package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserID(ctx *gin.Context) (uint, error) {
	userIDValue, exists := ctx.Get("userId")
	if !exists {
		return 0, errors.New("user not found in context")
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, errors.New("invalid user id type")
	}

	return userID, nil
}
