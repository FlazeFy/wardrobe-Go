package routes

import (
	"log"
	"wardrobe/models"

	"gorm.io/gorm"
)

func MigrateAll(db *gorm.DB) {
	db.AutoMigrate(
		&models.Admin{},
		&models.User{},
		&models.GoogleToken{},
		&models.Dictionary{},
		&models.Error{},
		&models.History{},
		&models.Feedback{},
		&models.Clothes{},
		&models.ClothesUsed{},
		&models.Outfit{},
		&models.OutfitRelation{},
		&models.OutfitUsed{},
		&models.Question{},
		&models.Schedule{},
		&models.UserRequest{},
		&models.UserTrack{},
		&models.UserWeather{},
		&models.Wash{},
		&models.AuditTrail{},
	)
	log.Printf("Migration : Success to migrate database")
}
