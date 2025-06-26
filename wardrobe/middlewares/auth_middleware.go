package middleware

import (
	"context"
	"net/http"
	"strings"
	"wardrobe/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(redisClient *redis.Client, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}

		// Check If Token Is Blacklisted In Redis
		val, err := redisClient.Get(context.Background(), tokenString).Result()
		if err == nil && val == "blacklisted" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token already expired"})
			return
		}

		// Validate JWT Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unexpected signing method"})
			}
			return config.GetJWTSecret(), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		// Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		// Extract userID
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user ID not found"})
			return
		}

		// Extract role
		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "role not found"})
			return
		}

		// Check If Role Is Allowed
		isAllowed := false
		for _, r := range allowedRoles {
			if role == r {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "access forbidden for this role"})
			return
		}

		// Set Context
		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}
