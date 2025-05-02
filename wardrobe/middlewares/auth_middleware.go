package middleware

import (
	"net/http"
	"strings"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization is required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		userId, err := utils.ValidateToken(parts[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
