package middleware

import (
	"log"
	"pelita/entity"
	"pelita/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuditTrailMiddleware(db *gorm.DB, activityName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Context User Id
		userID, err := utils.GetCurrentUserID(c)
		if err != nil {
			log.Println("Failed to get user ID from context:", err)
			c.Next()
			return
		}

		// Get Context Role
		userType, err := utils.GetCurrentRole(c)
		if err != nil {
			log.Println("Failed to get user role from context:", err)
			c.Next()
			return
		}

		history := entity.AuditTrail{
			ID:          uuid.New(),
			TypeUser:    userType,
			TypeAuditTrail: activityName,
			CreatedAt:   time.Now(),
		}

		// Fill ID Based Role
		switch userType {
		case "admin":
			history.AdminID = &userID
		case "user":
			history.UserID = &userID
		default:
			log.Println("unknown user type:")
			c.Next()
			return
		}

		err = db.Create(&history).Error

		if err != nil {
			log.Printf("failed to write audit log: %v\n", err)
		}

		c.Next()
	}
}
