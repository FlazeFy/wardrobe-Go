package utils

import (
	"time"
	"wardrobe/config"
	"wardrobe/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(config.GetJWTExpirationDuration()).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.GetJWTSecret())
}

func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := uint(claims["user_id"].(float64))
		return userId, nil
	}

	return 0, err
}

func HashPassword(u *models.User, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPass)
	return nil
}

func CheckPassword(u *models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
