package telegram

import (
	"log"
	"os"
	"wardrobe/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TelegramHandler() {
	// Init
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Callback Message
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			HandleStartCommand(update, bot)
		}
	}
}

func HandleStartCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userId := update.Message.Chat.ID
	firstName := update.Message.From.FirstName

	msg := tgbotapi.NewMessage(userId, "Hello "+firstName+"! Welcome to Wardrobe Bot")
	rows := [][]tgbotapi.InlineKeyboardButton{}

	// Menu List
	for _, item := range config.MenuList {
		btn := tgbotapi.NewInlineKeyboardButtonData(item.Label, item.Data)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	bot.Send(msg)
	bot.Send(tgbotapi.NewMessage(userId, "Choose an option : "))
}
