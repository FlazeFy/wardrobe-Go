package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserWeather struct {
		ID               uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
		WeatherTemp      float64   `json:"weather_temp" gorm:"type:float;not null"`
		WeatherHumid     int       `json:"weather_humid" gorm:"type:int;not null"`
		WeatherCity      string    `json:"weather_city" gorm:"type:varchar(75);not null"`
		WeatherCondition string    `json:"weather_condition" gorm:"type:varchar(16);not null"`
		WeatherHitFrom   string    `json:"weather_hit_from" gorm:"type:varchar(36);not null"`
		CreatedAt        time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		User      User      `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	// For Consume Open Weather API
	OpenWeatherAPIResponse struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
		City    string `json:"name"`
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`
	}
)
