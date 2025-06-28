package utils

import (
	"errors"
	"time"
	"wardrobe/config"
	"wardrobe/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId.String(),
		"exp":     time.Now().Add(config.GetJWTExpirationDuration()).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.GetJWTSecret())
}

func ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, errors.New("invalid token claims")
	}

	userIdStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.UUID{}, errors.New("user_id is not a string")
	}

	return uuid.Parse(userIdStr)
}

func HashPassword(user models.User, password string) (*models.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPass)

	return &user, nil
}

func CheckPassword(u *models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
