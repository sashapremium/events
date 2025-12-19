package contenteventsprocessor

import (
	"context"

	"github.com/sashapremium/events/internal/models"
)

type analyticsService interface {
	ProcessEvent(ctx context.Context, ev *models.ContentEvent) error
}
type ContentEventsProcessor struct {
	analyticsService analyticsService
}

func NewContentEventsProcessor(analyticsService analyticsService) *ContentEventsProcessor {
	return &ContentEventsProcessor{
		analyticsService: analyticsService,
	}
}
