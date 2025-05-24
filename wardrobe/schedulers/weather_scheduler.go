package schedulers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func BroadCastErrorToAdmin(db *gorm.DB) {
	// Get Admin Contact
	adminContext := utils.NewUserContext(db)
	contact, err := adminContext.GetAdminContact()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(contact) > 0 {
		for _, dt := range contact {
			if dt.TelegramUserId != nil && dt.TelegramIsValid {
				bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
				if err != nil {
					fmt.Println("Failed to connect to Telegram bot")
					return
				}

				telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
				if err != nil {
					fmt.Println("Invalid Telegram User Id")
					return
				}

				message := fmt.Sprintf("[ADMIN] Hello %s, there is an error in scheduler : weather_routine_fetch. Here's the detail :\n\n%s", dt.Username, err.Error())
				doc := tgbotapi.NewMessage(telegramID, message)
				doc.ParseMode = "html"

				_, err = bot.Send(doc)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}
	}
}

func SchedulerWeatherRoutineFetch() {
	db := config.ConnectDatabase()

	// Get User Ready Fetch Weather
	userContext := models.NewUserContext(db)
	users, err := userContext.SchedulerGetUserReadyFetchWeather()
	if err != nil {
		fmt.Println(err.Error())
		BroadCastErrorToAdmin(db)
		return
	}

	if len(users) > 0 {
		for _, dt := range users {
			open_weather_key := os.Getenv("OPEN_WEATHER_API_KEY")
			client := &http.Client{Timeout: 10 * time.Second}

			// Fetch OpenWeather Service
			url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&appid=%s", dt.TrackLat, dt.TrackLong, open_weather_key)
			resp, err := client.Get(url)
			if err != nil {
				fmt.Println("failed to request weather API: %w", err)
				BroadCastErrorToAdmin(db)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("failed to read weather API response: %w", err)
				BroadCastErrorToAdmin(db)
				return
			}
			var apiResp models.OpenWeatherAPIResponse
			if err := json.Unmarshal(body, &apiResp); err != nil {
				fmt.Println("failed to parse weather API JSON: %w", err)
				BroadCastErrorToAdmin(db)
				return
			}

			// Model Weather
			weather := &models.UserWeather{
				ID:               uuid.New(),
				WeatherTemp:      apiResp.Main.Temp,
				WeatherHumid:     apiResp.Main.Humidity,
				WeatherCity:      apiResp.City,
				WeatherCondition: apiResp.Weather[0].Main,
				WeatherHitFrom:   "Task Schedule",
				CreatedAt:        time.Now(),
				CreatedBy:        dt.UserID,
			}

			// Query : Create Weather
			if err := db.Create(&weather).Error; err != nil {
				fmt.Println("failed to create weather: %w", err)
				BroadCastErrorToAdmin(db)
				return
			}

			// Send to Telegram
			message := fmt.Sprintf("Hello %s, from your last coordinate %s,%s at %s. We've have checked the weather for today, and the result is:\n\nTemperature: %.2f °C\nHumidity: %d%s\nCity: %s\nWeather Condition: %s",
				dt.Username, dt.TrackLat, dt.TrackLong, dt.CreatedAt.Format("2006-01-02 15:04"), apiResp.Main.Temp, apiResp.Main.Humidity, "%", apiResp.City, apiResp.Weather[0].Main)
			if dt.TelegramUserId != nil && dt.TelegramIsValid {
				bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
				if err != nil {
					fmt.Println("Failed to connect to Telegram bot")
					BroadCastErrorToAdmin(db)
					return
				}

				telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
				if err != nil {
					fmt.Println("Invalid Telegram User Id")
					BroadCastErrorToAdmin(db)
					return
				}

				doc := tgbotapi.NewMessage(telegramID, message)
				doc.ParseMode = "html"

				_, err = bot.Send(doc)
				if err != nil {
					fmt.Println(err.Error())
					BroadCastErrorToAdmin(db)
					return
				}
			}
		}
	}
}
