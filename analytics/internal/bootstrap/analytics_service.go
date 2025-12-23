package bootstrap

import (
	"github.com/redis/go-redis/v9"
	analyticsService "github.com/sashapremium/events/analytics/internal/services/analyticsService"

	analyticsPG "github.com/sashapremium/events/analytics/internal/storage/analytics/pgstorage"
	"github.com/sashapremium/events/analytics/internal/storage/analytics/pgstorage/rediscache"
)

func InitAnalyticsService(storage *analyticsPG.PGStorage, redisClient *redis.Client) *analyticsService.Service {
	cache := rediscache.New(redisClient)
	return analyticsService.New(storage, cache)
}
