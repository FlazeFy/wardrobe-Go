package config

import (
	"os"
	"time"
)

func GetJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func GetJWTExpirationDuration() time.Duration {
	duration, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))

	if err != nil {
		return time.Hour * 24
	}

	return duration
}
