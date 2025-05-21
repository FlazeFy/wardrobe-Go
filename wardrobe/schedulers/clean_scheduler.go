package schedulers

import (
	"fmt"
	"os"
	"strconv"
	"wardrobe/config"
	"wardrobe/controllers"
	"wardrobe/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SchedulerCleanHistory() {
	days := 30
	db := config.ConnectDatabase()

	// Get User Contact
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
