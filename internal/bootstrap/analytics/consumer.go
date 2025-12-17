package analytics

import (
	"github.com/sashapremium/events/config"
	contenteventsconsumer "github.com/sashapremium/events/internal/consumer/content_events_consumer"
	contenteventsprocessor "github.com/sashapremium/events/internal/services/processors/content_events_processor"
)

func InitContentEventsConsumer(
	cfg *config.Config,
	processor *contenteventsprocessor.ContentEventsProcessor,
) *contenteventsconsumer.ContentEventsConsumer {
	return contenteventsconsumer.NewContentEventsConsumer(processor, cfg.Kafka.Brokers, cfg.Kafka.Topic)
}
