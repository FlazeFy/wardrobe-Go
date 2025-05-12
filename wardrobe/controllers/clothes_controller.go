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
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid clothes_price"})
			return
		}
		clothesPrice = &price
	}

	var clothesBuyAt *time.Time
	buyAtStr := ctx.PostForm("clothes_buy_at")
	if buyAtStr != "" {
		t, err := time.Parse(time.RFC3339, buyAtStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid clothes_buy_at"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "clothes_category must be one of: " + strings.Join(allowedCategories, ",")})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "clothes_type must be one of: " + strings.Join(allowedTypes, ",")})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "clothes_gender must be one of: " + strings.Join(allowedGenders, ",")})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "clothes_made_from must be one of: " + strings.Join(allowedMadeFroms, ",")})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "clothes_size must be one of: " + strings.Join(allowedSizes, ",")})
		return
	}

	// Query : Check Clothes Name
	var existing models.Clothes
	if err := c.DB.Where("clothes_name = ?", clothesName).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"message": "clothes with the same name already exists"})
		return
	}

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
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
			"message": "something went wrong",
		})
		return
	}

	// Get User Contact
	userContext := utils.NewUserContext(c.DB)
	contact, err := userContext.GetUserContact(*userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send to Telegram
	if contact.TelegramUserId != nil && contact.TelegramIsValid {
		filename := fmt.Sprintf("clothes-%s.pdf", clothes.ID)
		err = utils.GeneratePDF(clothes, filename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to connect to Telegram bot",
			})
			return
		}

		telegramID, err := strconv.ParseInt(*contact.TelegramUserId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Telegram User Id",
			})
			return
		}
		doc := tgbotapi.NewDocumentUpload(telegramID, filename)
		doc.Caption = fmt.Sprintf("clothes created, its called '%s'", clothes.ClothesName)

		_, err = bot.Send(doc)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send PDF to Telegram",
			})
			return
		}

		// Cleanup
		os.Remove(filename)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data":    clothes,
		"message": "clothes created",
	})
}

func (c *ClothesController) SoftDeleteClothesById(ctx *gin.Context) {
	// Param
	id := ctx.Param("id")

	// Models
	var clothes models.Clothes

	if err := c.DB.First(&clothes, "id = ?", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "clothes not found",
		})
		return
	}

	now := time.Now()
	clothes.DeletedAt = &now

	// Query : Update Clothes
	if err := c.DB.Save(&clothes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "clothes deleted",
	})
}

func (c *ClothesController) CreateClothesUsed(ctx *gin.Context) {
	// Models
	var req models.ClothesUsed

	// Validate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "used_context must be one of: " + strings.Join(allowedContexts, ",")})
		return
	}

	// Get User ID
	userId, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
			"message": "failed to create clothes used",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
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
	result := c.DB.Unscoped().First(&clothes_used, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "clothes used not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&clothes_used)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "clothes used permanentally deleted",
	})
}
