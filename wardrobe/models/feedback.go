package models

import "gorm.io/gorm"

type Feedback struct {
	gorm.Model
	FeedbackRate int    `gorm:"type:integer;not null"`
	FeedbackBody string `gorm:"type:varchar(255);not null"`
	// FK - User
	CreatedBy uint `json:"created_by" gorm:"not null"`
	User      User `json:"-" gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
