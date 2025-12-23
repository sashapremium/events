package bootstrap

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	addr := firstEnv("REDIS_ADDR", "ANALYTICS_REDIS_ADDR", "redis:6379")
	pass := firstEnv("REDIS_PASSWORD", "ANALYTICS_REDIS_PASSWORD", "")
	dbStr := firstEnv("REDIS_DB", "ANALYTICS_REDIS_DB", "0")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("ошибка подключения к Redis (%s): %v", addr, err)
	}

	return rdb
}

func firstEnv(primary, fallback, def string) string {
	if v := os.Getenv(primary); v != "" {
		return v
	}
	if v := os.Getenv(fallback); v != "" {
		return v
	}
	return def
}
