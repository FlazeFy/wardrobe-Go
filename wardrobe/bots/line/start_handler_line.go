package line

import (
	"net/http"
	"os"
	"wardrobe/config"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func init() {
	var err error
	// Load Env
	err = godotenv.Load()
	if err != nil {
		panic("Error loading ENV")
	}

	// Init
	bot, err = linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		panic(err)
	}
}

func LineHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		var username string

		// Callback Message
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				if event.Source.Type == linebot.EventSourceTypeUser {
					profile, _ := bot.GetProfile(event.Source.UserID).Do()
					username = profile.DisplayName
				} else {
					username = ""
				}

				if message, ok := event.Message.(*linebot.TextMessage); ok {
					if message.Text == "/start" {
						carousel := makeCarousel(username)
						_, err := bot.ReplyMessage(event.ReplyToken, carousel).Do()
						if err != nil {
							c.String(http.StatusInternalServerError, err.Error())
						}
					}
				}
			}
		}

		utils.BuildResponseMessage(c, "success", "line bot", "get", http.StatusOK, nil, nil)
	}
}

func makeCarousel(username string) *linebot.TemplateMessage {
	columns := []*linebot.CarouselColumn{}
	chunkSize := 2

	for i := 0; i < len(config.MenuList); i += chunkSize {
		end := i + chunkSize
		if end > len(config.MenuList) {
			end = len(config.MenuList)
		}
		chunk := config.MenuList[i:end]
		actions := []linebot.TemplateAction{}

		// Menu List
		for _, item := range chunk {
			actions = append(actions, linebot.NewMessageAction(item.Label, item.Data))
		}

		columns = append(columns, &linebot.CarouselColumn{
			Title:   "Hello " + username + "! Welcome to Wardrobe Bot",
			Text:    "Choose an option : ",
			Actions: actions,
		})
	}

	return linebot.NewTemplateMessage("Select a menu", &linebot.CarouselTemplate{Columns: columns})
}
