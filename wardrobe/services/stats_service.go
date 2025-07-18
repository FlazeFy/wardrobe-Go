package services

import (
	"encoding/json"
	"wardrobe/cache"
	"wardrobe/models/others"
	"wardrobe/repositories"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Stats Interface
type StatsService interface {
	GetMostUsedContext(tableName, targetCol string, userId uuid.UUID) ([]others.StatsContextTotal, error)
	GetMonthlyClothesUsedByClothesIdAndYear(year int, tableName, targetCol, contextCol, contextId string, userId uuid.UUID) ([]others.StatsContextTotal, error)
}

// Stats Struct
type statsService struct {
	statsRepo   repositories.StatsRepository
	redisClient *redis.Client
	statsCache  cache.StatsCache
}

// Stats Constructor
func NewStatsService(statsRepo repositories.StatsRepository, redisClient *redis.Client, statsCache cache.StatsCache) StatsService {
	return &statsService{
		statsRepo:   statsRepo,
		redisClient: redisClient,
		statsCache:  statsCache,
	}
}

func (s *statsService) GetMostUsedContext(tableName, targetCol string, userId uuid.UUID) ([]others.StatsContextTotal, error) {
	// Cache : Get Key
	cacheKey := s.statsCache.StatsKeyMostUsedContext("clothes", targetCol, userId)
	// Cache : Temp Stats
	stats, err := s.statsCache.GetStatsMostUsedContext(s.redisClient, cacheKey)
	if err == nil {
		return stats, nil
	}

	// Repo : Find Most Used Context
	stats, err = s.statsRepo.FindMostUsedContext(tableName, targetCol, userId)
	if err != nil {
		return nil, err
	}

	// Cache : Store Redis
	jsonData, _ := json.Marshal(stats)
	s.statsCache.SetStatsMostUsedContext(s.redisClient, cacheKey, jsonData)

	return stats, nil
}

func (s *statsService) GetMonthlyClothesUsedByClothesIdAndYear(year int, tableName, targetCol, contextCol, contextId string, userId uuid.UUID) ([]others.StatsContextTotal, error) {
	// Cache : Get Key
	cacheKey := s.statsCache.StatsKeyMostUsedContext("monthly"+"_"+tableName, string(year)+"_"+contextCol+"_"+contextId, userId)
	// Cache : Temp Stats
	stats, err := s.statsCache.GetStatsMostUsedContext(s.redisClient, cacheKey)
	if err == nil {
		return stats, nil
	}

	// Repo : Find Monthly Clothes Used By Clothes Id And Year
	stats, err = s.statsRepo.FindMonthlyClothesUsedByClothesIdAndYear(year, tableName, targetCol, contextCol, contextId, userId)
	if err != nil {
		return nil, err
	}

	// Cache : Store Redis
	jsonData, _ := json.Marshal(stats)
	s.statsCache.SetStatsMostUsedContext(s.redisClient, cacheKey, jsonData)

	return stats, nil
}
