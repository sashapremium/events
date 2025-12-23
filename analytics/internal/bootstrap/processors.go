package bootstrap

import (
	analyticsService "github.com/sashapremium/events/analytics/internal/services/analyticsService"
	processor "github.com/sashapremium/events/analytics/internal/services/processors/content_events_processor"
)

func InitEventsProcessor(svc *analyticsService.Service) *processor.ContentEventsProcessor {
	return processor.NewContentEventsProcessor(svc)
}
