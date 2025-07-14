package repositories

import (
	"fmt"
	"wardrobe/models/others"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Stats Interface
type StatsRepository interface {
	FindMostUsedContext(tableName, targetCol string, userId uuid.UUID) ([]others.StatsContextTotal, error)
}

// Stats Struct
type statsRepository struct {
	db *gorm.DB
}

// Stats Constructor
func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) FindMostUsedContext(tableName, targetCol string, userId uuid.UUID) ([]others.StatsContextTotal, error) {
	// Models
	var stats []others.StatsContextTotal

	// Query
	result := r.db.Table(tableName).
		Select(fmt.Sprintf("COUNT(%s) as total, %s as context", targetCol, targetCol)).
		Where("created_by", userId).
		Group(targetCol).
		Order("total DESC").
		Limit(7).
		Find(&stats)

	if result.Error != nil {
		return nil, result.Error
	}
	if len(stats) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return stats, nil
}
