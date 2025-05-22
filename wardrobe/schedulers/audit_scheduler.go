package schedulers

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SchedulerAuditError() {
	db := config.ConnectDatabase()

	// Get Admin Contact
	userContext := utils.NewUserContext(db)
	contact, err := userContext.GetAdminContact()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Get All Error
	errorContext := models.NewErrorContext(db)
	errors, err := errorContext.GetAllErrorAudit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Send to Telegram
	datetime := time.Now()
	if len(contact) > 0 && len(errors) > 0 {
		filename := fmt.Sprintf("audit_error_%s.pdf", datetime)
		err = utils.GeneratePDFErrorAudit(errors, filename)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

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

				doc := tgbotapi.NewDocumentUpload(telegramID, filename)
				doc.Caption = fmt.Sprintf("[ADMIN] Hello %s, the system just run an audit error, with result of %d error found. Here's the document", dt.Username, len(errors))

				_, err = bot.Send(doc)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}

		// Cleanup
		os.Remove(filename)
	}
}
