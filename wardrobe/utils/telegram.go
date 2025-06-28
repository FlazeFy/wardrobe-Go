package utils

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendTelegramTextMessage(telegramUserID, message string, typeMessage string, file *interface{}) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Printf("Failed to connect to Telegram bot: %v", err)
	}

	telegramID, err := strconv.ParseInt(telegramUserID, 10, 64)
	if err != nil {
		log.Printf("Invalid Telegram User Id")
	}

	switch typeMessage {
	case "text":
		doc := tgbotapi.NewMessage(telegramID, message)
		_, err = bot.Send(doc)
	case "doc":
		doc := tgbotapi.NewDocumentUpload(telegramID, file)
		doc.Caption = message
		_, err = bot.Send(doc)
	}

	if err != nil {
		log.Printf("Failed to send message to Telegram: %v", err)
	}
}
