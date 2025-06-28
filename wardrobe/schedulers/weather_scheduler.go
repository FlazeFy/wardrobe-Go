package schedulers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"wardrobe/models"
	"wardrobe/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type WeatherScheduler struct {
	AdminService       services.AdminService
	UserService        services.UserService
	UserWeatherService services.UserWeatherService
}

func NewWeatherScheduler(
	adminService services.AdminService,
	userService services.UserService,
	userWeatherService services.UserWeatherService,
) *WeatherScheduler {
	return &WeatherScheduler{
		AdminService:       adminService,
		UserService:        userService,
		UserWeatherService: userWeatherService,
	}
}

func (s *WeatherScheduler) BroadCastErrorToAdmin() {
	// Service : Get All Admin Contact
	contact, err := s.AdminService.GetAllAdminContact()
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

func (s *WeatherScheduler) SchedulerWeatherRoutineFetch() {
	// Get User Ready Fetch Weather
	users, err := s.UserService.SchedulerGetUserReadyFetchWeather()
	if err != nil {
		fmt.Println(err.Error())
		s.BroadCastErrorToAdmin()
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
				s.BroadCastErrorToAdmin()
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("failed to read weather API response: %w", err)
				s.BroadCastErrorToAdmin()
				return
			}
			var apiResp models.OpenWeatherAPIResponse
			if err := json.Unmarshal(body, &apiResp); err != nil {
				fmt.Println("failed to parse weather API JSON: %w", err)
				s.BroadCastErrorToAdmin()
				return
			}

			// Model Weather
			weather := &models.UserWeather{
				WeatherTemp:      apiResp.Main.Temp,
				WeatherHumid:     apiResp.Main.Humidity,
				WeatherCity:      apiResp.City,
				WeatherCondition: apiResp.Weather[0].Main,
				WeatherHitFrom:   "Task Schedule",
			}

			// Query : Create Weather
			err = s.UserWeatherService.Create(weather, dt.UserID)
			if err != nil {
				fmt.Println("failed to create weather: %w", err)
				s.BroadCastErrorToAdmin()
				return
			}

			// Send to Telegram
			message := fmt.Sprintf("Hello %s, from your last coordinate %s,%s at %s. We've have checked the weather for today, and the result is:\n\nTemperature: %.2f °C\nHumidity: %d%s\nCity: %s\nWeather Condition: %s",
				dt.Username, dt.TrackLat, dt.TrackLong, dt.CreatedAt.Format("2006-01-02 15:04"), apiResp.Main.Temp, apiResp.Main.Humidity, "%", apiResp.City, apiResp.Weather[0].Main)
			if dt.TelegramUserId != nil && dt.TelegramIsValid {
				bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
				if err != nil {
					fmt.Println("Failed to connect to Telegram bot")
					s.BroadCastErrorToAdmin()
					return
				}

				telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
				if err != nil {
					fmt.Println("Invalid Telegram User Id")
					s.BroadCastErrorToAdmin()
					return
				}

				doc := tgbotapi.NewMessage(telegramID, message)
				doc.ParseMode = "html"

				_, err = bot.Send(doc)
				if err != nil {
					fmt.Println(err.Error())
					s.BroadCastErrorToAdmin()
					return
				}
			}
		}
	}
}
