package main

import (
	"log"
	"os"

	"github.com/sashapremium/events/config"
	"github.com/sashapremium/events/internal/bootstrap/events"
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
	// Хранилище
	eventsStorage := events.InitPGStorage(cfg)
	// Kafka-продюсер
	eventsProducer := events.InitKafkaProducer(cfg)
	//Сервис
	eventsService := events.InitEventsService(eventsStorage, eventsProducer)
	//API
	eventsAPI := events.InitEventsServiceAPI(eventsService)

	if err := events.AppRun(cfg, eventsAPI); err != nil {
		log.Fatal(err)
	}

}
