package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	AuditTrail struct {
		ID           uuid.UUID  `json:"id" gorm:"type:varchar(36);primaryKey"`
		AdminID      *uuid.UUID `json:"admin_id" gorm:"type:varchar(36);null"`
		UserID       *uuid.UUID `json:"user_id" gorm:"type:varchar(36);null"`
		TypeUser     string     `json:"type_user" gorm:"type:varchar(36);not null"`
		TypeHistory  string     `json:"type_history" gorm:"type:varchar(255);not null"`
		CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	}
)
