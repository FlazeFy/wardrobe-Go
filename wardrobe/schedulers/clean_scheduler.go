package schedulers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"wardrobe/config"
	"wardrobe/controllers"
	"wardrobe/models"
	"wardrobe/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SchedulerCleanHistory() {
	days := 30
	db := config.ConnectDatabase()

	// Get Admin Contact
	userContext := utils.NewUserContext(db)
	contact, err := userContext.GetAdminContact()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Delete History
	historyContext := controllers.NewHistoryController(db)
	total, err := historyContext.DeleteHistoryForLastNDays(days)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Send to Telegram
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

				msgText := fmt.Sprintf("[ADMIN] Hello %s, the system just run a clean history, with result of %d history executed", dt.Username, total)
				msg := tgbotapi.NewMessage(telegramID, msgText)

				_, err = bot.Send(msg)
				if err != nil {
					fmt.Println("Failed to send message to Telegram")
					return
				}
			}
		}
	}
}

func SchedulerCleanDeletedClothes() {
	days := 30
	db := config.ConnectDatabase()

	// Get User Contact
	userContext := utils.NewUserContext(db)
	contact, err := userContext.GetAdminContact()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Get Clothes Plan Destroy
	clothesContext := models.NewClothesContext(db)
	plans, err := clothesContext.GetClothesPlanDestroy(days)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Send to Telegram
	if len(plans) > 0 {
		totalUser := 0
		userBefore := ""
		listClothes := ""
		totalClothes := 0
		var next *models.ClothesPlanDestroy

		scheduleContext := models.NewScheduleContext(db)
		washContext := models.NewWashContext(db)

		for idx, dt := range plans {
			_, err = clothesContext.SchedulerHardDeleteClothesById(dt.ID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			_, err = clothesContext.SchedulerDeleteClothesUsedByClothesId(dt.ID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			_, err = scheduleContext.SchedulerDeleteSchedulehById(dt.ID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			_, err = washContext.SchedulerDeleteWashById(dt.ID)
			if err != nil {
				fmt.Println(err.Error())
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
						fmt.Println("Failed to connect to Telegram bot")
						return
					}

					telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
					if err != nil {
						fmt.Println("Invalid Telegram User Id")
						return
					}

					msgText := fmt.Sprintf("Hello %s, We've recently cleaned up your deleted clothes. Here are the details:\n\n%s", dt.Username, listClothes)
					msg := tgbotapi.NewMessage(telegramID, msgText)

					_, err = bot.Send(msg)
					if err != nil {
						fmt.Println("Failed to send message to Telegram")
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
						fmt.Println("Failed to connect to Telegram bot")
						return
					}

					telegramID, err := strconv.ParseInt(*dt.TelegramUserId, 10, 64)
					if err != nil {
						fmt.Println("Invalid Telegram User Id")
						return
					}

					msgText := fmt.Sprintf("[ADMIN] Hello %s, the system just run a clean deleted clothes, with result of %d clothes executed from %d user", dt.Username, totalClothes, totalUser)
					msg := tgbotapi.NewMessage(telegramID, msgText)

					_, err = bot.Send(msg)
					if err != nil {
						fmt.Println("Failed to send message to Telegram")
						return
					}
				}
			}
		}
	}
}
