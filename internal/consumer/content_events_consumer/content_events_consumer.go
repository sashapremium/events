package contenteventsconsumer

import (
	"context"

	"github.com/sashapremium/events/internal/models"
)

type contentEventsProcessor interface {
	Handle(ctx context.Context, ev *models.ContentEvent) error
}

type ContentEventsConsumer struct {
	processor   contentEventsProcessor
	kafkaBroker []string
	topicName   string
	groupID     string
}

func NewContentEventsConsumer(
	processor contentEventsProcessor,
	kafkaBroker []string,
	topicName string,
) *ContentEventsConsumer {
	return &ContentEventsConsumer{
		processor:   processor,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
		groupID:     "AnalyticsService_group",
	}
}
