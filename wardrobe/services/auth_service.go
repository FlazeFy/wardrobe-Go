package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/models/others"
	"wardrobe/repositories"
	"wardrobe/utils"

	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService interface {
	BasicRegister(userReq models.User) (*string, error)
	BasicSignOut(token string) error
	BasicLogin(loginReq others.LoginRequest) (*string, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	adminRepo   repositories.AdminRepository
	redisClient *redis.Client
}

func NewAuthService(userRepo repositories.UserRepository, adminRepo repositories.AdminRepository, redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:    userRepo,
		adminRepo:   adminRepo,
		redisClient: redisClient,
	}
}

func (s *authService) BasicRegister(userReq models.User) (*string, error) {
	// Repo : Find By Email
	user, err := s.userRepo.FindByUsernameOrEmail(userReq.Username, userReq.Email)
	if user != nil || err != gorm.ErrRecordNotFound {
		if user != nil {
			return nil, errors.New("username or email has already been used")
		}

		return nil, err
	}

	// Hashing
	user, err = utils.HashPassword(userReq, userReq.Password)
	if err != nil {
		return nil, errors.New("failed hashing password")
	}

	// Service : Create User
	user, err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// JWT Token Generate
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed generating token")
	}

	return &token, nil
}

func (s *authService) BasicLogin(loginReq others.LoginRequest) (*string, error) {
	// Repo : Find By Email
	user, err := s.userRepo.FindByEmail(loginReq.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("account not found")
		}

		return nil, err
	}

	// Utils : Check Password
	if err := utils.CheckPassword(user, loginReq.Password); err != nil {
		return nil, errors.New("invalid password")
	}

	// Utils : JWT Token Generate
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed generating token")
	}

	return &token, nil
}

func (s *authService) BasicSignOut(authHeader string) error {
	// Clean Bearer
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return errors.New("invalid authorization header")
	}

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
