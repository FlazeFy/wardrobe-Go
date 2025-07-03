package middleware

import (
	"log"
	"time"
	"wardrobe/models"
	"wardrobe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuditTrailMiddleware(db *gorm.DB, activityName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Context User Id
		userID, err := utils.GetUserID(c)
		if err != nil {
			log.Println("Failed to get user ID from context:", err)
			c.Next()
			return
		}

		auditTrail := models.AuditTrail{
			ID:             uuid.New(),
			UserID:         userID,
			TypeAuditTrail: activityName,
			CreatedAt:      time.Now(),
		}

		err = db.Create(&auditTrail).Error
		if err != nil {
			log.Printf("failed to write audit log: %v\n", err)
		}

		c.Next()
	}
}
