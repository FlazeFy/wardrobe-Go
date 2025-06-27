package services

import (
	"errors"
	"fmt"
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schedule Interface
type ScheduleService interface {
	CreateSchedule(schedule models.Schedule, userID uuid.UUID) error
	GetScheduleByDay(day string, userId uuid.UUID) ([]models.ScheduleByDay, error)
	DeleteScheduleByClothesId(id uuid.UUID) (int64, error)
	HardDeleteScheduleById(id, createdBy uuid.UUID) error
}

// Schedule Struct
type scheduleService struct {
	scheduleRepo repositories.ScheduleRepository
	userRepo     repositories.UserRepository
	clothesRepo  repositories.ClothesRepository
}

// Schedule Constructor
func NewScheduleService(scheduleRepo repositories.ScheduleRepository, userRepo repositories.UserRepository, clothesRepo repositories.ClothesRepository) ScheduleService {
	return &scheduleService{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
		clothesRepo:  clothesRepo,
	}
}

func (s *scheduleService) GetScheduleByDay(day string, userId uuid.UUID) ([]models.ScheduleByDay, error) {
	return s.scheduleRepo.FindScheduleByDay(day, userId)
}

func (s *scheduleService) CreateSchedule(req models.Schedule, userID uuid.UUID) error {
	// Repo : Check Schedule By Day And Clothes ID
	isFound, err := s.scheduleRepo.CheckScheduleByDayAndClothesID(req.Day, userID, req.ClothesId)
	if err != nil {
		return err
	}
	if isFound {
		return errors.New("schedule with same day already exist")
	}

	// Repo : Create Schedule
	err = s.scheduleRepo.CreateSchedule(&req, userID)
	if err != nil {
		return err
	}

	// Repo : Find User Contact By Id
	contact, err := s.userRepo.FindUserContactByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user contact not found")
		}

		return err
	}

	// Repo : Find Clothes Short Info By Id
	clothes, err := s.clothesRepo.FindClothesShortInfoById(req.ClothesId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("clothes not found")
		}

		return err
	}

	// Utils : Send to Telegram
	if contact.TelegramUserId != nil && contact.TelegramIsValid {
		message := fmt.Sprintf("Your clothes called '%s' has been added to weekly schedule and set to wear on every %s", clothes.ClothesName, req.Day)
		utils.SendTelegramTextMessage(*contact.TelegramUserId, message)
	}

	return nil
}

func (s *scheduleService) DeleteScheduleByClothesId(id uuid.UUID) (int64, error) {
	// Repo : Delete Schedule By Clothes Id
	rows, err := s.scheduleRepo.DeleteScheduleByClothesId(id)
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *scheduleService) HardDeleteScheduleById(id, createdBy uuid.UUID) error {
	return s.scheduleRepo.HardDeleteScheduleById(id, createdBy)
}
