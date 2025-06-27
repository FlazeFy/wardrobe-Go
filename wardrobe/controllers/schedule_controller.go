package controllers

import (
	"errors"
	"net/http"
	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/services"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleController struct {
	ScheduleService services.ScheduleService
}

func NewScheduleController(scheduleService services.ScheduleService) *ScheduleController {
	return &ScheduleController{ScheduleService: scheduleService}
}

// Query
func (c *ScheduleController) GetScheduleByDay(ctx *gin.Context) {
	// Params
	day := ctx.Param("day")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "feedback", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Get Schedule By Day
	data, err := c.ScheduleService.GetScheduleByDay(day, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "schedule", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "schedule", "get", http.StatusOK, data, nil)
}

func (c *ScheduleController) GetScheduleForTomorrow(ctx *gin.Context) {
	// Params
	day := ctx.Param("day")

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "feedback", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Helper : Get Next Days
	tomorrow := utils.GetNextDay(day, 1)
	twoDaysLater := utils.GetNextDay(day, 2)

	// Service : Get Schedule By Day
	dataTommorow, errTomorrow := c.ScheduleService.GetScheduleByDay(tomorrow, *userID)
	if errTomorrow != nil && !errors.Is(errTomorrow, gorm.ErrRecordNotFound) {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}
	dataTwoDayLater, errTwoDayLater := c.ScheduleService.GetScheduleByDay(tomorrow, *userID)
	if errTwoDayLater != nil && !errors.Is(errTwoDayLater, gorm.ErrRecordNotFound) {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(ctx, "success", "schedule", "get", http.StatusOK, gin.H{
		"tomorrow":           utils.CheckIfEmpty(dataTommorow),
		"tomorrow_day":       tomorrow,
		"two_days_later":     utils.CheckIfEmpty(dataTwoDayLater),
		"two_days_later_day": twoDaysLater,
	}, nil)
}

// Command
func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {
	// Models
	var req models.Schedule

	// Validate JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BuildResponseMessage(ctx, "failed", "feedback", "invalid request body", http.StatusBadRequest, nil, nil)
		return
	}

	// Validate Field
	if req.Day == "" {
		utils.BuildResponseMessage(ctx, "failed", "schedule", "day is required", http.StatusBadRequest, nil, nil)
		return
	}
	if req.ClothesId == uuid.Nil {
		utils.BuildResponseMessage(ctx, "failed", "schedule", "clothes_id is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Validator Contain : Day
	if !utils.Contains(config.DictionaryTypes, req.Day) {
		utils.BuildResponseMessage(ctx, "failed", "schedule", "day is not valid", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "feedback", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Create Schedule
	err = c.ScheduleService.CreateSchedule(req, *userID)
	if err != nil {
		if err.Error() == "schedule with same day already exist" {
			utils.BuildResponseMessage(ctx, "failed", "schedule", err.Error(), http.StatusConflict, nil, nil)
			return
		}
		if err.Error() == "user contact not found" || err.Error() == "clothes not found" {
			utils.BuildResponseMessage(ctx, "failed", "schedule", err.Error(), http.StatusNotFound, nil, nil)
			return
		}

		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "schedule", "post", http.StatusCreated, nil, nil)
}

func (c *ScheduleController) HardDeleteScheduleById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Parse Param UUID
	scheduleID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "schedule", "invalid id", http.StatusBadRequest, nil, nil)
		return
	}

	// Get User ID
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.BuildResponseMessage(ctx, "failed", "schedule", err.Error(), http.StatusBadRequest, nil, nil)
		return
	}

	// Service : Hard Delete History By ID
	err = c.ScheduleService.HardDeleteScheduleById(scheduleID, *userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildResponseMessage(ctx, "failed", "schedule", "empty", http.StatusNotFound, nil, nil)
		return
	}
	if err != nil {
		utils.BuildErrorMessage(ctx, err.Error())
		return
	}

	utils.BuildResponseMessage(ctx, "success", "schedule", "hard delete", http.StatusOK, nil, nil)
}
