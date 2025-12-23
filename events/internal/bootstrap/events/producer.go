package events

import (
	"github.com/sashapremium/events/events/config"
	"github.com/sashapremium/events/events/internal/services/processors/content_events_processor"
)

func InitKafkaProducer(cfg *config.Config) contenteventsprocessor.Producer {
	return contenteventsprocessor.NewKafkaProducer(cfg.Kafka)
}
