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

type ClothesController struct {
	DB *gorm.DB
}

func NewClothesController(db *gorm.DB) *ClothesController {
	return &ClothesController{DB: db}
}

// Query
func (c *ClothesController) GetClothesLastHistory(ctx *gin.Context) {
	// Get User ID
	userIdStr, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	userId := *userIdStr

	// Query : Get Last Created
	clothesContext := models.NewClothesContext(c.DB)
	resLastAdded, err := clothesContext.GetClothesLastCreated("created_at", userId)
	if err != nil {
		if err.Error() != "clothes not found" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "something went wrong",
			})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	// Query : Get Last Deleted
	resLastDeleted, _ := clothesContext.GetClothesLastDeleted("deleted_at", userId)

	// Response
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

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes last history fetched",
		"data":    data,
	})
}

func (c *ClothesController) GetDeletedClothes(ctx *gin.Context) {
	// Get User ID
	userIdStr, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	userId := *userIdStr

	// Query : Get Deleted Clothes
	clothesContext := models.NewClothesContext(c.DB)
	res, err := clothesContext.GetDeletedClothes(userId)

	// Response
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes fetched",
		"data":    res,
	})
}

func (c *ClothesController) GetAllClothesHeader(ctx *gin.Context) {
	// Param
	category := ctx.Param("category")
	order := ctx.Param("order")

	// Get User ID
	userIdStr, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	userId := *userIdStr

	// Query : Get All Clothes Header
	clothesContext := models.NewClothesContext(c.DB)
	res, err := clothesContext.GetAllClothesHeader(category, order, userId)
	if err != nil {
		if err.Error() != "clothes not found" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "something went wrong",
			})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes fetched",
		"data":    res,
	})
}

func (c *ClothesController) GetAllClothesDetail(ctx *gin.Context) {
	// Param
	category := ctx.Param("category")
	order := ctx.Param("order")

	// Get User ID
	userIdStr, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	userId := *userIdStr

	// Query : Get All Clothes Detail
	clothesContext := models.NewClothesContext(c.DB)
	res, err := clothesContext.GetAllClothesDetail(category, order, userId)
	if err != nil {
		if err.Error() != "clothes not found" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "something went wrong",
			})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes fetched",
		"data":    res,
	})
}

// Command
func (c *ClothesController) CreateClothes(ctx *gin.Context) {
	// Mandatory Field
	clothesName := ctx.PostForm("clothes_name")
	clothesColor := ctx.PostForm("clothes_color")
	clothesMadeFrom := ctx.PostForm("clothes_made_from")
	clothesType := ctx.PostForm("clothes_type")
	clothesCategory := ctx.PostForm("clothes_category")
	clothesSize := ctx.PostForm("clothes_size")
	clothesGender := ctx.PostForm("clothes_gender")
	clothesQty, _ := strconv.Atoi(ctx.PostForm("clothes_qty"))
	isFaded, _ := strconv.ParseBool(ctx.PostForm("is_faded"))
	hasWashed, _ := strconv.ParseBool(ctx.PostForm("has_washed"))
	hasIroned, _ := strconv.ParseBool(ctx.PostForm("has_ironed"))
	isFavorite, _ := strconv.ParseBool(ctx.PostForm("is_favorite"))
	isScheduled, _ := strconv.ParseBool(ctx.PostForm("is_scheduled"))

	// Optional Field
	var clothesDesc *string
	desc := ctx.PostForm("clothes_desc")
	if desc != "" {
		clothesDesc = &desc
	}

	var clothesMerk *string
	merk := ctx.PostForm("clothes_merk")
	if merk != "" {
		clothesMerk = &merk
	}

	var clothesPrice *int
	priceStr := ctx.PostForm("clothes_price")
	if priceStr != "" {
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "invalid clothes_price",
			})
			return
		}
		clothesPrice = &price
	}

	var clothesBuyAt *time.Time
	buyAtStr := ctx.PostForm("clothes_buy_at")
	if buyAtStr != "" {
		t, err := time.Parse(time.RFC3339, buyAtStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "invalid clothes_buy_at",
			})
			return
		}
		clothesBuyAt = &t
	}

	// Validate : Clothes Category Rules
	allowedCategories := []string{"upper_body", "bottom_body", "head", "foot", "hand"}
	ok := false
	for _, v := range allowedCategories {
		if clothesCategory == v {
			ok = true
			break
		}
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "clothes_category must be one of: " + strings.Join(allowedCategories, ","),
		})
		return
	}

	// Validate : Clothes Type Rules
	allowedTypes := []string{"hat", "pants", "shirt", "jacket", "shoes", "socks", "scarf", "gloves", "shorts", "skirt", "dress", "blouse", "sweater", "hoodie", "tie", "belt", "coat", "underwear", "swimsuit", "vest", "t-shirt", "jeans", "leggings", "boots", "sandals", "sneakers", "raincoat", "poncho", "cardigan"}
	ok = false
	for _, v := range allowedTypes {
		if clothesType == v {
			ok = true
			break
		}
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "clothes_type must be one of: " + strings.Join(allowedTypes, ","),
		})
		return
	}

	// Validate : Clothes Gender Rules
	allowedGenders := []string{"male", "female", "unisex"}
	ok = false
	for _, v := range allowedGenders {
		if clothesGender == v {
			ok = true
			break
		}
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "clothes_gender must be one of: " + strings.Join(allowedGenders, ","),
		})
		return
	}

	// Validate : Clothes Made From Rules
	allowedMadeFroms := []string{"cotton", "wool", "silk", "linen", "polyester", "denim", "leather", "nylon", "rayon", "synthetic", "cloth"}
	ok = false
	for _, v := range allowedMadeFroms {
		if clothesMadeFrom == v {
			ok = true
			break
		}
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "clothes_made_from must be one of: " + strings.Join(allowedMadeFroms, ","),
		})
		return
	}

	// Validate : Clothes Size Rules
	allowedSizes := []string{"S", "M", "L", "XL", "XXL"}
	ok = false
	for _, v := range allowedSizes {
		if clothesSize == v {
			ok = true
			break
		}
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "clothes_size must be one of: " + strings.Join(allowedSizes, ","),
		})
		return
	}

	// Query : Check Clothes Name
	var existing models.Clothes
	if err := c.DB.Where("clothes_name = ?", clothesName).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"status":  "failed",
			"message": "clothes with the same name already exists",
		})
		return
	}

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	clothes := models.Clothes{
		ID:              uuid.New(),
		ClothesName:     clothesName,
		ClothesDesc:     clothesDesc,
		ClothesMerk:     clothesMerk,
		ClothesColor:    clothesColor,
		ClothesPrice:    clothesPrice,
		ClothesBuyAt:    clothesBuyAt,
		ClothesQty:      clothesQty,
		ClothesImage:    nil,
		IsFaded:         isFaded,
		HasWashed:       hasWashed,
		HasIroned:       hasIroned,
		IsFavorite:      isFavorite,
		IsScheduled:     isScheduled,
		CreatedAt:       time.Now(),
		UpdatedAt:       nil,
		DeletedAt:       nil,
		CreatedBy:       *userId,
		ClothesMadeFrom: clothesMadeFrom,
		ClothesType:     clothesType,
		ClothesCategory: clothesCategory,
		ClothesSize:     clothesSize,
		ClothesGender:   clothesGender,
	}

	// Query : Create Clothes
	if err := c.DB.Create(&clothes).Error; err != nil {
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

	// Send to Telegram
	if contact.TelegramUserId != nil && contact.TelegramIsValid {
		filename := fmt.Sprintf("clothes-%s.pdf", clothes.ID)
		err = utils.GeneratePDF(clothes, filename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

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
		doc := tgbotapi.NewDocumentUpload(telegramID, filename)
		doc.Caption = fmt.Sprintf("clothes created, its called '%s'", clothes.ClothesName)

		_, err = bot.Send(doc)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Failed to send PDF to Telegram",
			})
			return
		}

		// Cleanup
		os.Remove(filename)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    clothes,
		"message": "clothes created",
	})
}

func (c *ClothesController) SoftDeleteClothesById(ctx *gin.Context) {
	// Param
	id := ctx.Param("id")

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	// Models
	var clothes models.Clothes

	if err := c.DB.First(&clothes, "id = ? AND deleted_at is null AND created_by = ?", id, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "clothes not found",
		})
		return
	}

	now := time.Now()
	clothes.DeletedAt = &now

	// Query : Update Clothes
	if err := c.DB.Save(&clothes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "something went wrong",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes deleted",
	})
}

func (c *ClothesController) RecoverDeletedClothesById(ctx *gin.Context) {
	// Param
	id := ctx.Param("id")

	// Models
	var clothes models.Clothes

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	if err := c.DB.First(&clothes, "id = ? AND deleted_at is not null AND created_by = ?", id, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "clothes not found",
		})
		return
	}

	clothes.DeletedAt = nil

	// Query : Update Clothes
	if err := c.DB.Save(&clothes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "something went wrong",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes recovered",
	})
}

func (c *ClothesController) HardDeleteClothesById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid user id",
		})
		return
	}

	// Models
	var clothes models.Clothes
	var schedule models.Schedule
	var clothes_used models.ClothesUsed
	var wash models.Wash
	var outfit_rel models.OutfitRelation

	uuidID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid UUID",
		})
		return
	}

	// Get Clothes
	clothesContext := models.NewClothesContext(c.DB)
	clothes_old, err := clothesContext.GetClothesShortInfoById(uuidID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Query : Delete Clothes
	result := c.DB.Unscoped().First(&clothes, "id = ? AND deleted_at is not null AND created_by = ?", id, userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "clothes not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&clothes)

	// Query : Delete Clothes Relation
	c.DB.Unscoped().Where("clothes_id = ? AND created_by = ?", id, userId).Delete(&schedule)
	c.DB.Unscoped().Where("clothes_id = ? AND created_by = ?", id, userId).Delete(&clothes_used)
	c.DB.Unscoped().Where("clothes_id = ? AND created_by = ?", id, userId).Delete(&wash)
	c.DB.Unscoped().Where("clothes_id = ? AND created_by = ?", id, userId).Delete(&outfit_rel)

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
		message := fmt.Sprintf("Your clothes called '%s' has been permentally removed from Wardrobe", clothes_old.ClothesName)
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

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes permanentally deleted",
	})
}

func (c *ClothesController) CreateClothesUsed(ctx *gin.Context) {
	// Models
	var req models.ClothesUsed

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request body",
		})
		return
	}

	// Validate : Clothes Category Rules
	allowedContexts := []string{"Worship", "Shopping", "Work", "School", "Campus", "Sport", "Party"}
	ok := false
	for _, v := range allowedContexts {
		if req.UsedContext == v {
			ok = true
			break
		}
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "used_context must be one of: " + strings.Join(allowedContexts, ","),
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

	// Query : Add Clothes Used
	clothes_used := models.ClothesUsed{
		ID:          uuid.New(),
		ClothesNote: req.ClothesNote,
		ClothesId:   req.ClothesId,
		UsedContext: req.UsedContext,
		CreatedBy:   *userId,
	}
	if err := c.DB.Create(&clothes_used).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to create clothes used",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    clothes_used,
		"message": "clothes used created",
	})
}

func (c *ClothesController) HardDeleteClothesUsedById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Models
	var clothes_used models.ClothesUsed

	// Query
	result := c.DB.Unscoped().First(&clothes_used, "id = ? AND deleted_at is not null", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "clothes used not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&clothes_used)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes used permanentally deleted",
	})
}

func (c *ClothesController) GetClothesUsedHistory(ctx *gin.Context) {
	// Params
	clothes_id_param := ctx.Param("clothes_id")
	order := ctx.Param("order")

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	// Clothes Id
	var clothes_id uuid.UUID
	if clothes_id_param == "all" {
		clothes_id = uuid.Nil
	} else {
		clothes_id, err = uuid.Parse(clothes_id_param)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "invalid clothes id",
			})
			return
		}
	}

	// Query : Get Clothes Used History
	clothesContext := models.NewClothesUsedContext(c.DB)
	res, err := clothesContext.GetClothesUsedHistory(*userId, clothes_id, order)

	// Response
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "clothes fetched",
		"data":    res,
	})
}
