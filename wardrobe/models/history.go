package models

import "gorm.io/gorm"

type History struct {
	gorm.Model
	HistoryType    string `gorm:"type:varchar(36);not null"`
	HistoryContext string `gorm:"type:varchar(255);not null"`
	// FK - User
	CreatedBy uint `json:"created_by" gorm:"not null"`
	User      User `json:"user" gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
