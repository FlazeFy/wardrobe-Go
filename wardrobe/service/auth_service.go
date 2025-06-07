package service

import (
	"context"
	"errors"
	"time"

	"wardrobe/config"
	"wardrobe/entity"
	"wardrobe/repository"
	"wardrobe/utils"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	SignOut(token string) error
}

type authService struct {
	userRepo       repository.UserRepository
	adminRepo      repository.AdminRepository
	redisClient    *redis.Client
}

func NewAuthService(userRepo repository.UserRepository, adminRepo repository.AdminRepository,redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:       userRepo,
		adminRepo:      adminRepo,
		redisClient:    redisClient,
	}
}

func (s *authService) SignOut(tokenString string) error {
	// Token Parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("missing exp in token")
	}

	// Check If Token Expired
	expTime := time.Unix(int64(expFloat), 0)
	duration := time.Until(expTime)
	if duration <= 0 {
		return errors.New("token already expired")
	}

	// Save token to Redis blacklist
	err = s.redisClient.Set(context.Background(), tokenString, "blacklisted", duration).Err()
	if err != nil {
		return errors.New("failed to blacklist token")
	}

	return nil
}
