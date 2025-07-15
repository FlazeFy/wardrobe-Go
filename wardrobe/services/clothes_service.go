package services

import (
	"errors"
	"fmt"
	"os"
	"time"
	"wardrobe/models"
	"wardrobe/repositories"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Clothes Interface
type ClothesService interface {
	CreateClothes(clothes *models.Clothes, userID uuid.UUID) (*models.Clothes, error)
	GetDeletedClothes(userID uuid.UUID) ([]models.ClothesDeleted, error)
	GetAllClothesHeader(pagination utils.Pagination, category, order string, userID uuid.UUID) ([]models.ClothesHeader, int64, error)
	GetAllClothesDetail(pagination utils.Pagination, category, order string, userID uuid.UUID) ([]models.Clothes, int64, error)
	GetClothesLastCreated(ctx string, userID uuid.UUID) (*models.ClothesLastCreated, error)
	GetClothesLastHistory(userID uuid.UUID) (interface{}, error)
	SoftDeleteClothesById(userID, clothesID uuid.UUID) error
	RecoverDeletedClothesById(userID, clothesID uuid.UUID) error
	HardDeleteClothesById(userID, clothesID uuid.UUID) error

	// For Task Scheduler
	GetClothesPlanDestroy(days int) ([]models.ClothesPlanDestroy, error)
	SchedulerHardDeleteClothesById(id uuid.UUID) (int64, error)
	SchedulerGetUnusedClothes(days int) ([]models.SchedulerClothesUnused, error)
	SchedulerGetUnironedClothes() ([]models.SchedulerClothesUnironed, error)
}

// Clothes Struct
type clothesService struct {
	clothesRepo        repositories.ClothesRepository
	userRepo           repositories.UserRepository
	scheduleRepo       repositories.ScheduleRepository
	clothesUsedRepo    repositories.ClothesUsedRepository
	washRepo           repositories.WashRepository
	outfitRelationRepo repositories.OutfitRelationRepository
}

// Clothes Constructor
func NewClothesService(clothesRepo repositories.ClothesRepository, userRepo repositories.UserRepository, scheduleRepo repositories.ScheduleRepository,
	clothesUsedRepo repositories.ClothesUsedRepository, washRepo repositories.WashRepository, outfitRelationRepo repositories.OutfitRelationRepository) ClothesService {
	return &clothesService{
		clothesRepo:        clothesRepo,
		userRepo:           userRepo,
		scheduleRepo:       scheduleRepo,
		clothesUsedRepo:    clothesUsedRepo,
		washRepo:           washRepo,
		outfitRelationRepo: outfitRelationRepo,
	}
}

func (s *clothesService) GetAllClothesHeader(pagination utils.Pagination, category, order string, userID uuid.UUID) ([]models.ClothesHeader, int64, error) {
	return s.clothesRepo.FindAllClothesHeader(pagination, category, order, userID)
}

func (s *clothesService) GetAllClothesDetail(pagination utils.Pagination, category, order string, userID uuid.UUID) ([]models.Clothes, int64, error) {
	return s.clothesRepo.FindAllClothesDetail(pagination, category, order, userID)
}

func (s *clothesService) GetClothesLastCreated(ctx string, userID uuid.UUID) (*models.ClothesLastCreated, error) {
	return s.clothesRepo.FindClothesLastCreated(ctx, userID)
}

func (s *clothesService) GetDeletedClothes(userID uuid.UUID) ([]models.ClothesDeleted, error) {
	return s.clothesRepo.FindDeletedClothes(userID)
}

func (s *clothesService) CreateClothes(clothes *models.Clothes, userID uuid.UUID) (*models.Clothes, error) {
	// Repo : Check Clothes By Name
	isFound, err := s.clothesRepo.CheckClothesByName(clothes.ClothesName, userID)
	if err != nil {
		return nil, err
	}
	if isFound {
		return nil, errors.New("clothes with the same name already exists")
	}

	// Repo : Create Clothes
	clothes, err = s.clothesRepo.CreateClothes(clothes, userID)
	if err != nil {
		return nil, err
	}

	// Repo : Find User Contact By Id
	contact, err := s.userRepo.FindUserContactByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user contact not found")
		}

		return nil, err
	}

	// Send to Telegram
	if contact.TelegramUserId != nil && contact.TelegramIsValid {
		filename := fmt.Sprintf("clothes-%s.pdf", clothes.ID)
		err = utils.GeneratePDFCreateClothes(clothes, filename)
		if err != nil {
			return nil, err
		}

		message := fmt.Sprintf("clothes created, its called '%s'", clothes.ClothesName)
		utils.SendTelegramTextMessage(*contact.TelegramUserId, message, "doc", nil)

		// Cleanup
		os.Remove(filename)
	}

	return clothes, err
}

func (s *clothesService) HardDeleteClothesById(clothesID, userID uuid.UUID) error {
	// Get Clothes
	clothes_old, err := s.clothesRepo.FindClothesShortInfoById(clothesID)
	if err != nil {
		return err
	}

	// Service : Hard Delete Clothes By Id
	err = s.clothesRepo.HardDeleteClothesById(clothesID, userID)
	if err != nil {
		return err
	}

	// Service : Hard Delete Schedule By Clothes ID
	err = s.scheduleRepo.HardDeleteScheduleByClothesID(clothesID, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// Service : Hard Delete Clothes Used By Clothes ID
	err = s.clothesUsedRepo.HardDeleteClothesUsedByClothesID(clothesID, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// Service : Hard Delete Wash By Clothes ID
	err = s.washRepo.HardDeleteWashByClothesID(clothesID, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// Service : Hard Delete Outfit Relation By Clothes ID
	err = s.outfitRelationRepo.HardDeleteOutfitRelationByClothesID(clothesID, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
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

	// Send to Telegram
	if contact.TelegramUserId != nil && contact.TelegramIsValid {
		message := fmt.Sprintf("Your clothes called '%s' has been permentally removed from Wardrobe", clothes_old.ClothesName)
		utils.SendTelegramTextMessage(*contact.TelegramUserId, message, "text", nil)
	}

	return err
}

func (s *clothesService) SoftDeleteClothesById(userID, clothesID uuid.UUID) error {
	// Repo : Find Clothes By Id
	clothes, err := s.clothesRepo.FindClothesById(clothesID)
	if err != nil {
		return err
	}
	if clothes.DeletedAt != nil {
		return gorm.ErrRecordNotFound
	}

	// Soft Delete
	now := time.Now()
	clothes.DeletedAt = &now

	// Repo : Update Clothes By Id
	err = s.clothesRepo.UpdateClothesById(clothes, clothesID)
	if err != nil {
		return err
	}

	return nil
}

func (s *clothesService) RecoverDeletedClothesById(userID, clothesID uuid.UUID) error {
	// Repo : Find Clothes By Id
	clothes, err := s.clothesRepo.FindClothesById(clothesID)
	if err != nil {
		return err
	}
	if clothes.DeletedAt == nil {
		return gorm.ErrRecordNotFound
	}

	// Recover
	clothes.DeletedAt = nil

	// Repo : Update Clothes By Id
	err = s.clothesRepo.UpdateClothesById(clothes, clothesID)
	if err != nil {
		return err
	}

	return nil
}

func (s *clothesService) GetClothesLastHistory(userID uuid.UUID) (interface{}, error) {
	// Repo : Find Last Created
	resLastAdded, err := s.clothesRepo.FindClothesLastCreated("created_at", userID)
	if err != nil {
		return nil, err
	}

	// Repo : Find Last Deleted
	resLastDeleted, err := s.clothesRepo.FindClothesLastDeleted("deleted_at", userID)
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	data := gin.H{
		"last_added_clothes":   resLastAdded.ClothesName,
		"last_added_date":      resLastAdded.CreatedAt,
		"last_deleted_clothes": nil,
		"last_deleted_date":    nil,
	}

	if resLastDeleted != nil {
		data["last_deleted_clothes"] = resLastDeleted.ClothesName
		data["last_deleted_date"] = resLastDeleted.DeletedAt
	}

	return data, nil
}

// For Task Scheduler
func (s *clothesService) GetClothesPlanDestroy(days int) ([]models.ClothesPlanDestroy, error) {
	// Repo : Find Clothes Plan Destroy
	rows, err := s.clothesRepo.FindClothesPlanDestroy(days)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *clothesService) SchedulerHardDeleteClothesById(id uuid.UUID) (int64, error) {
	// Repo : Hard Delete Clothes By Id
	rows, err := s.clothesRepo.HardDeleteClothesById2(id)
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *clothesService) SchedulerGetUnusedClothes(days int) ([]models.SchedulerClothesUnused, error) {
	// Repo : Scheduler Find Unused Clothes
	rows, err := s.clothesRepo.SchedulerFindUnusedClothes(days)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *clothesService) SchedulerGetUnironedClothes() ([]models.SchedulerClothesUnironed, error) {
	// Repo : Scheduler Find Unironed Clothes
	rows, err := s.clothesRepo.SchedulerFindUnironedClothes()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
