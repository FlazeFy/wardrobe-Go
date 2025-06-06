package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleController struct {
	DB *gorm.DB
}

func NewScheduleController(db *gorm.DB) *ScheduleController {
	return &ScheduleController{DB: db}
}

// Query
func (c *ScheduleController) GetScheduleByDay(ctx *gin.Context) {
	// Params
	day := ctx.Param("day")

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Query
	scheduleContext := models.NewScheduleContext(c.DB)
	res, err := scheduleContext.GetScheduleByDay(day, *userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Response
	if len(res) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "schedule not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    res,
		"message": "schedule fetched",
	})
}

func (c *ScheduleController) GetScheduleForTomorrow(ctx *gin.Context) {
	// Params
	day := ctx.Param("day")

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Helper : Get Next Days
	tomorrow := utils.GetNextDay(day, 1)
	twoDaysLater := utils.GetNextDay(day, 2)

	// Query : Get Schedule
	scheduleContext := models.NewScheduleContext(c.DB)
	resTomorrow, err := scheduleContext.GetScheduleByDay(tomorrow, *userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	resTwoDaysLater, err := scheduleContext.GetScheduleByDay(twoDaysLater, *userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"tomorrow":           utils.CheckIfEmpty(resTomorrow),
			"tomorrow_day":       tomorrow,
			"two_days_later":     utils.CheckIfEmpty(resTwoDaysLater),
			"two_days_later_day": twoDaysLater,
		},
	})
}

// Command
func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {
	// Models
	var req models.Schedule

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request body",
		})
		return
	}

	day := req.Day
	// Validate : Day Rules
	allowedDays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	ok := false
	for _, v := range allowedDays {
		if day == v {
			ok = true
			break
		}
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "day must be one of: " + strings.Join(allowedDays, ","),
		})
		return
	}

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Query : Check Schedule
	clothes_id := req.ClothesId
	var existing models.Schedule
	if err := c.DB.Where("day = ? AND created_by = ? AND clothes_id = ?", day, userId, clothes_id).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"message": "schedule with the same name already exists"})
		return
	}

	schedule := models.Schedule{
		ID:           uuid.New(),
		ClothesId:    clothes_id,
		Day:          day,
		ScheduleNote: req.ScheduleNote,
		IsRemind:     req.IsRemind,
		CreatedAt:    time.Now(),
		CreatedBy:    *userId,
	}

	// Query : Create Schedule
	if err := c.DB.Create(&schedule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "something went wrong",
		})
		return
	}

	// Get User Contact
	userContext := utils.NewUserContext(c.DB)
	contact, err := userContext.GetUserContact(*userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Get Clothes
	clothesContext := models.NewClothesContext(c.DB)
	clothes, err := clothesContext.GetClothesShortInfoById(clothes_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Send to Telegram
	if contact.TelegramUserId != nil && contact.TelegramIsValid {
		bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Failed to connect to Telegram bot",
			})
			return
		}

		telegramID, err := strconv.ParseInt(*contact.TelegramUserId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "Invalid Telegram User Id",
			})
			return
		}
		message := fmt.Sprintf("Your clothes called '%s' has been added to weekly schedule and set to wear on every %s", clothes.ClothesName, day)
		doc := tgbotapi.NewMessage(telegramID, message)

		_, err = bot.Send(doc)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Failed to send message to Telegram",
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    schedule,
		"message": "schedule created",
	})
}

func (c *ScheduleController) HardDeleteScheduleById(ctx *gin.Context) {
	// Param
	id := ctx.Param("id")

	// Models
	var schedule models.Schedule

	// Query : Delete Schedule
	result := c.DB.Unscoped().First(&schedule, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "schedule not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&schedule)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "schedule permentally delete",
	})
}
