package schedulers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"wardrobe/services"
	"wardrobe/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CalendarScheduler struct {
	AdminService    services.AdminService
	ScheduleService services.ScheduleService
}

func NewCalendarScheduler(
	adminService services.AdminService,
	scheduleService services.ScheduleService,
) *CalendarScheduler {
	return &CalendarScheduler{
		AdminService:    adminService,
		ScheduleService: scheduleService,
	}
}

func (s *CalendarScheduler) SchedulerCalendarSycnSchedule() {
	// Service : Get All Admin Contact
	contact, err := s.AdminService.GetAllAdminContact()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Service : Find Schedule Ready To Assign Calendar Task By Day
	day := utils.GetTomorrowDayName()
	schedules, err := s.ScheduleService.FindScheduleReadyToAssignCalendarTaskByDay(day)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Utils : Sync the Google Calendar
	var list_sync = ""
	var total_failed = 0
	var total_success = 0

	if len(schedules) > 0 {
		for _, dt := range schedules {
			// Prepare Task Calendar Props
			summary := fmt.Sprintf("- Schedule To Wear : %s at %s\n", dt.ClothesName, dt.Day)
			desc := fmt.Sprintf("Notes : %s", *dt.ScheduleNote)
			startTime, err := utils.GetThisWeekdayWithHour(dt.Day, 8, 0)
			if err != nil {
				log.Printf("Invalid day format:", dt.Day)
				total_failed++
				continue
			}

			// Sync To Calendar
			_, err = utils.AddWeeklyGoogleCalendarEvent(dt.AccessToken, summary, desc, startTime)
			if err != nil {
				log.Printf("Failed to add calendar event:", err)
				total_failed++
			} else {
				// Service : Update Remind By Id
				err := s.ScheduleService.UpdateRemindByID(dt.ID, true)
				if err != nil {
					log.Printf("Failed to add calendar event:", err)
					total_failed++
				} else {
					total_success++
					list_sync += fmt.Sprintf("- %s for %s\n", dt.Username, dt.Day)
				}
			}
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

					msgText := fmt.Sprintf("[ADMIN] Hello %s, We're here to inform you. We have successfully added sync Schedule with User's Google Calendar Account. Here's the detail :\n\n%s", dt.Username, list_sync)
					doc := tgbotapi.NewMessage(telegramID, msgText)
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
}
