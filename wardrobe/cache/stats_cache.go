package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"wardrobe/config"
	"wardrobe/models/others"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Stats Interface
type StatsCache interface {
	StatsKeyMostUsedContext(table, field string, userID uuid.UUID) string
	RefreshStatsMostUsedContext(redis *redis.Client, table string, fields []string, userID uuid.UUID)
	SetStatsMostUsedContext(redis *redis.Client, key string, data []byte) error
	GetStatsMostUsedContext(redis *redis.Client, key string) ([]others.StatsContextTotal, error)
}

// Stats Struct
type statsCache struct {
	redisClient *redis.Client
}

// Stats Constructor
func NewStatsCache(redisClient *redis.Client) StatsCache {
	return &statsCache{
		redisClient: redisClient,
	}
}

// For API GET : Most Used Context
func (ch *statsCache) StatsKeyMostUsedContext(table, field string, userID uuid.UUID) string {
	return fmt.Sprintf("stats:%s:%s:%s", table, field, userID.String())
}

func (ch *statsCache) RefreshStatsMostUsedContext(redis *redis.Client, table string, fields []string, userID uuid.UUID) {
	for _, field := range fields {
		key := ch.StatsKeyMostUsedContext(table, field, userID)
		redis.Del(context.Background(), key)
	}
}

func (ch *statsCache) SetStatsMostUsedContext(redis *redis.Client, key string, data []byte) error {
	return redis.Set(context.Background(), key, data, config.RedisTime).Err()
}

func (ch *statsCache) GetStatsMostUsedContext(redis *redis.Client, key string) ([]others.StatsContextTotal, error) {
	var stats []others.StatsContextTotal

	cached, err := redis.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(cached), &stats); err != nil {
		return nil, err
	}

	return stats, nil
}
