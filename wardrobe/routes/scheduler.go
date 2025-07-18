package routes

import (
	"time"
	"wardrobe/schedulers"
	"wardrobe/services"

	"github.com/robfig/cron"
)

func SetUpScheduler(
	adminService services.AdminService, errorService services.ErrorService, historyService services.HistoryService, clothesService services.ClothesService,
	clothesUsedService services.ClothesUsedService, scheduleService services.ScheduleService, washService services.WashService, questionService services.QuestionService,
	userService services.UserService, userWeatherService services.UserWeatherService) {

	// Initialize Scheduler
	auditScheduler := schedulers.NewAuditScheduler(errorService, adminService)
	cleanScheduler := schedulers.NewCleanScheduler(adminService, historyService, clothesService, clothesUsedService, scheduleService, washService)
	reminderScheduler := schedulers.NewReminderScheduler(adminService, clothesService, clothesUsedService, questionService)
	weatherScheduler := schedulers.NewWeatherScheduler(adminService, userService, userWeatherService)
	calendarScheduler := schedulers.NewCalendarScheduler(adminService, scheduleService)
	houseKeepingScheduler := schedulers.NewHouseKeepingScheduler(adminService)

	// Init Scheduler
	c := cron.New()
	Scheduler(c, auditScheduler, cleanScheduler, reminderScheduler, weatherScheduler, calendarScheduler, houseKeepingScheduler)
	c.Start()
	defer c.Stop()
}

func Scheduler(c *cron.Cron, auditScheduler *schedulers.AuditScheduler, cleanScheduler *schedulers.CleanScheduler, reminderScheduler *schedulers.ReminderScheduler,
	weatherScheduler *schedulers.WeatherScheduler, calendarScheduler *schedulers.CalendarScheduler, houseKeepingScheduler *schedulers.HouseKeepingScheduler) {
	// For Production
	// Clean Scheduler
	c.AddFunc("10 0 * * *", weatherScheduler.SchedulerWeatherRoutineFetch)
	c.AddFunc("0 2 * * *", cleanScheduler.SchedulerCleanHistory)
	c.AddFunc("0 2 * * *", cleanScheduler.SchedulerCleanDeletedClothes)
	c.AddFunc("0 1 * * 1", auditScheduler.SchedulerAuditError)
	c.AddFunc("20 2 * * 1,3,6", reminderScheduler.SchedulerReminderUnansweredQuestion)
	c.AddFunc("20 2 * * 0,2,5", reminderScheduler.SchedulerReminderUnusedClothes)
	c.AddFunc("20 3 * * *", reminderScheduler.SchedulerReminderUnironedClothes)
	c.AddFunc("0 3 * * *", reminderScheduler.SchedulerReminderWashUsedClothes)
	c.AddFunc("45 1 * * *", calendarScheduler.SchedulerCalendarSycnSchedule)
	c.AddFunc("0 5 2 * *", houseKeepingScheduler.SchedulerMonthlyLog)

	// For Development
	go func() {
		time.Sleep(5 * time.Second)

		// Clean Scheduler
		// cleanScheduler.SchedulerCleanHistory()
		// cleanScheduler.SchedulerCleanDeletedClothes()

		// Audit Scheduler
		// auditScheduler.SchedulerAuditError()

		// Reminder Scheduler
		// reminderScheduler.SchedulerReminderUnansweredQuestion()
		// reminderScheduler.SchedulerReminderUnusedClothes()
		// reminderScheduler.SchedulerReminderUnironedClothes()
		// reminderScheduler.SchedulerReminderWashUsedClothes()

		// Weather Scheduler
		// weatherScheduler.SchedulerWeatherRoutineFetch()

		// Calendar Scheduler
		// calendarScheduler.SchedulerCalendarSycnSchedule()

		// House Keeping Scheduler
		houseKeepingScheduler.SchedulerMonthlyLog()
	}()
}
