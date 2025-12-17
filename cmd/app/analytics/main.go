package analytics

import (
	"log"
	"os"

	"github.com/sashapremium/events/config"
	"github.com/sashapremium/events/internal/bootstrap/analytics"
)

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func main() {
	cfgPath := getenv("CONFIG_PATH", "config.yaml")

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("ошибка парсинга конфига %q: %v", cfgPath, err)
	}

	pg := analytics.InitPGStorage(cfg)
	rdb := analytics.InitRedis()

	svc := analytics.InitAnalyticsService(pg, rdb)
	processor := analytics.InitEventsProcessor(svc)
	consumer := analytics.InitContentEventsConsumer(cfg, processor)
	api := analytics.InitAnalyticsServiceAPI(svc)

	if err := analytics.AppRun(cfg, api, consumer); err != nil {
		log.Fatal(err)
	}
}
