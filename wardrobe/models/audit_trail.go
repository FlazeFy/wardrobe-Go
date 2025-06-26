package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	AuditTrail struct {
		ID             uuid.UUID  `json:"id" gorm:"type:varchar(36);primaryKey"`
		UserID         *uuid.UUID `json:"user_id" gorm:"type:varchar(36);null"`
		TypeAuditTrail string     `json:"type_history" gorm:"type:varchar(255);not null"`
		CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	}
)
