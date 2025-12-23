package main

import (
	"log"
	"os"

	"github.com/sashapremium/events/analytics/config"
	"github.com/sashapremium/events/analytics/internal/bootstrap"
)

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func main() {
	cfgPath := getenv("CONFIG_PATH", "/config/analytics.yaml")

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("ошибка парсинга конфига %q: %v", cfgPath, err)
	}

	pg := bootstrap.InitPGStorage(cfg)
	rdb := bootstrap.InitRedis()
	svc := bootstrap.InitAnalyticsService(pg, rdb)
	processor := bootstrap.InitEventsProcessor(svc)
	consumer := bootstrap.InitContentEventsConsumer(cfg, processor)
	api := bootstrap.InitAnalyticsServiceAPI(svc)
	bootstrap.AppRun(api, consumer)
}
