package bootstrap

import (
	"github.com/sashapremium/events/analytics/config"
	contenteventsconsumer "github.com/sashapremium/events/analytics/internal/consumer/content_events_consumer"
	contenteventsprocessor "github.com/sashapremium/events/analytics/internal/services/processors/content_events_processor"
)

func InitContentEventsConsumer(
	cfg *config.Config,
	processor *contenteventsprocessor.ContentEventsProcessor,
) *contenteventsconsumer.ContentEventsConsumer {
	return contenteventsconsumer.NewContentEventsConsumer(processor, cfg.Kafka.Brokers, cfg.Kafka.Topic)
}
