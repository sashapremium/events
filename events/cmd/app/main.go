package main

import (
	"log"
	"os"

	"github.com/sashapremium/events/events/config"
	"github.com/sashapremium/events/events/internal/bootstrap"
)

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func main() {
	cfgPath := getenv("CONFIG_PATH", "/config/config.yaml")

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("ошибка парсинга конфига %q: %v", cfgPath, err)
	}
	// Хранилище
	eventsStorage := bootstrap.InitPGStorage(cfg)
	// Kafka-продюсер
	eventsProducer := bootstrap.InitKafkaProducer(cfg)
	//Сервис
	eventsService := bootstrap.InitEventsService(eventsStorage, eventsProducer)
	//API
	eventsAPI := bootstrap.InitEventsServiceAPI(eventsService)

	bootstrap.AppRun(eventsAPI)

}
