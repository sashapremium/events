package contenteventsprocessor

import (
	"context"

	"github.com/sashapremium/events/events/internal/models"
)

type analyticsService interface {
	ProcessEvent(ctx context.Context, ev *models.ContentEvent) error
}
type ContentEventsProcessor struct {
	analyticsService analyticsService
}
