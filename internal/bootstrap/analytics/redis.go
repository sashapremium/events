package analytics

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	addr := getenv(
		"ANALYTICS_REDIS_ADDR", "localhost:6379")
	pass := getenv("ANALYTICS_REDIS_PASSWORD", "")
	dbStr := getenv("ANALYTICS_REDIS_DB", "0")

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
		log.Fatalf("ошибка подключения к Redis: %v", err)
	}

	return rdb
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
