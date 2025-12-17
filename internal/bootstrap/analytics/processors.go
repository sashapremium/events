package analytics

import (
	analyticsService "github.com/sashapremium/events/internal/services/analyticsService"
	processor "github.com/sashapremium/events/internal/services/processors/content_events_processor"
)

func InitEventsProcessor(svc *analyticsService.Service) *processor.ContentEventsProcessor {
	return processor.NewContentEventsProcessor(svc)
}
