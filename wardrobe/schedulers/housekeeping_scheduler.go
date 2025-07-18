package schedulers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"wardrobe/services"
	"wardrobe/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HouseKeepingScheduler struct {
	AdminService services.AdminService
}

func NewHouseKeepingScheduler(
	adminService services.AdminService,
) *HouseKeepingScheduler {
	return &HouseKeepingScheduler{
		AdminService: adminService,
	}
}

func (s *HouseKeepingScheduler) SchedulerMonthlyLog() {
	// Service : Get All Admin Contact
	contact, err := s.AdminService.GetAllAdminContact()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Helpers : Clean Logs
	logPath, err := utils.GetLastMonthLogFilePath()
	if err != nil {
		log.Println("Log file not found:", err)
		return
	}

	// Open the log file
	fileBytes, err := os.Open(logPath)
	if err != nil {
		log.Println("Failed to open log file:", err)
		return
	}
	defer fileBytes.Close()

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

				file, err := os.Open(logPath)
				if err != nil {
					log.Println("Failed to open log file:", err)
					return
				}
				defer file.Close()

				fileInfo, err := file.Stat()
				if err != nil {
					log.Println("Failed to stat log file:", err)
					return
				}

				fileReader := tgbotapi.FileReader{
					Name:   fileInfo.Name(),
					Reader: file,
					Size:   fileInfo.Size(),
				}

				doc := tgbotapi.NewDocumentUpload(telegramID, fileReader)
				doc.ParseMode = "html"
				doc.Caption = fmt.Sprintf("[ADMIN] Hello %s, here is housekeeping log for %s %d",
					dt.Username, time.Now().AddDate(0, -1, 0).Format("January"), time.Now().AddDate(0, -1, 0).Year())

				_, err = bot.Send(doc)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
		}

		if err := utils.DeleteFileByPath(logPath); err != nil {
			log.Println("Failed to delete log file:", err)
		}
	}
}
