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

func SchedulerReminderUnansweredQuestion() {
	db := config.ConnectDatabase()

	// Get Admin Contact
	userContext := utils.NewUserContext(db)
	contact, err := userContext.GetAdminContact()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Get Unanswered Question
	questionContext := models.NewQuestionContext(db)
	questions, err := questionContext.GetUnansweredQuestion()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Send to Telegram
	datetime := time.Now()
	if len(contact) > 0 && len(questions) > 0 {
		filename := fmt.Sprintf("reminder_unanswered_question_%s.pdf", datetime)
		err = utils.GeneratePDFReminderUnansweredQuestion(questions, filename)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var list_question = ""
		for _, dt := range questions {
			list_question += fmt.Sprintf("- %s\nNotes: <i>ask at %s</i>\n\n", dt.Question, dt.CreatedAt)
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
				doc.ParseMode = "html"
				doc.Caption = fmt.Sprintf("[ADMIN] Hello %s, We're here to remind you. You have some unanswered question that needed to be answer. Here are the details:\n\n%s", dt.Username, list_question)

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
