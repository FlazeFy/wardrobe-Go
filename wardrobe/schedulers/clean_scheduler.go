package schedulers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"wardrobe/models"
	"wardrobe/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CleanScheduler struct {
	AdminService       services.AdminService
	HistoryService     services.HistoryService
	ClothesService     services.ClothesService
	ClothesUsedService services.ClothesUsedService
	ScheduleService    services.ScheduleService
	WashService        services.WashService
}

func NewCleanScheduler(
	adminService services.AdminService,
	historyService services.HistoryService,
	clothesService services.ClothesService,
	clothesUsedService services.ClothesUsedService,
	scheduleService services.ScheduleService,
	washService services.WashService,
) *CleanScheduler {
	return &CleanScheduler{
		AdminService:       adminService,
		HistoryService:     historyService,
		ClothesService:     clothesService,
		ClothesUsedService: clothesUsedService,
		ScheduleService:    scheduleService,
		WashService:        washService,
	}
}

func (s *CleanScheduler) SchedulerCleanHistory() {
	days := 30

	// Service : Get All Admin Contact
	contact, err := s.AdminService.GetAllAdminContact()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Service : Delete History For Last N Days
	total, err := s.HistoryService.DeleteHistoryForLastNDays(days)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Send to Telegram
	if len(contact) > 0 {
		for _, dt := range contact {
			if dt.TelegramUserId != nil && dt.TelegramIsValid {
				bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
				if err != nil {
					log.Println("Failed to connect to Telegram bot")
					return
				}

				telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
				if err != nil {
					log.Println("Invalid Telegram User Id")
					return
				}

				msgText := fmt.Sprintf("[ADMIN] Hello %s, the system just run a clean history, with result of %d history executed", dt.Username, total)
				msg := tgbotapi.NewMessage(telegramID, msgText)

				_, err = bot.Send(msg)
				if err != nil {
					log.Println("Failed to send message to Telegram")
					return
				}
			}
		}
	}
}

func (s *CleanScheduler) SchedulerCleanDeletedClothes() {
	days := 30

	// Service : Get All Admin Contact
	contact, err := s.AdminService.GetAllAdminContact()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Service : Get Clothes Plan Destroy
	plans, err := s.ClothesService.GetClothesPlanDestroy(days)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Send to Telegram
	if len(plans) > 0 {
		totalUser := 0
		userBefore := ""
		listClothes := ""
		totalClothes := 0
		var next *models.ClothesPlanDestroy

		for idx, dt := range plans {
			// Service : Scheduler Hard Delete Clothes By Id
			_, err = s.ClothesService.SchedulerHardDeleteClothesById(dt.ID)
			if err != nil {
				log.Println(err.Error())
				return
			}

			// Service : Scheduler Hard Delete Clothes Used By Clothes Id
			_, err = s.ClothesUsedService.DeleteClothesUsedByClothesId(dt.ID)
			if err != nil {
				log.Println(err.Error())
				return
			}

			// Service : Scheduler Hard Delete Schedule By Clothes Id
			_, err = s.ScheduleService.DeleteScheduleByClothesId(dt.ID)
			if err != nil {
				log.Println(err.Error())
				return
			}

			// Service : Scheduler Hard Delete Wash By Clothes Id
			_, err = s.WashService.DeleteWashByClothesId(dt.ID)
			if err != nil {
				log.Println(err.Error())
				return
			}
			totalClothes++

			if userBefore == "" || userBefore == dt.Username {
				listClothes += fmt.Sprintf("- %s\n", strings.Title(dt.ClothesName))
			}

			if idx+1 < len(plans) {
				next = &plans[idx+1]
			} else {
				next = nil
			}

			isLastOrDiffUser := next == nil || next.Username != dt.Username

			if isLastOrDiffUser {
				if dt.TelegramUserId != nil && dt.TelegramIsValid {
					bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
					if err != nil {
						log.Println("Failed to connect to Telegram bot")
						return
					}

					telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
					if err != nil {
						log.Println("Invalid Telegram User Id")
						return
					}

					msgText := fmt.Sprintf("Hello %s, We've recently cleaned up your deleted clothes. Here are the details:\n\n%s", dt.Username, listClothes)
					msg := tgbotapi.NewMessage(telegramID, msgText)

					_, err = bot.Send(msg)
					if err != nil {
						log.Println("Failed to send message to Telegram")
						return
					}
					totalUser++
				}

				listClothes = ""
				totalClothes++
			}

			userBefore = dt.Username
		}

		if len(contact) > 0 {
			for _, dt := range contact {
				if dt.TelegramUserId != nil && dt.TelegramIsValid {
					bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
					if err != nil {
						log.Println("Failed to connect to Telegram bot")
						return
					}

					telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
					if err != nil {
						log.Println("Invalid Telegram User Id")
						return
					}

					msgText := fmt.Sprintf("[ADMIN] Hello %s, the system just run a clean deleted clothes, with result of %d clothes executed from %d user", dt.Username, totalClothes, totalUser)
					msg := tgbotapi.NewMessage(telegramID, msgText)

					_, err = bot.Send(msg)
					if err != nil {
						log.Println("Failed to send message to Telegram")
						return
					}
				}
			}
		}
	}
}
