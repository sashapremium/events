package contenteventsprocessor

import (
	"context"

	"github.com/sashapremium/events/analytics/internal/models"
)

func (p *ContentEventsProcessor) Handle(ctx context.Context, ev *models.ContentEvent) error {
	return p.analyticsService.ProcessEvent(ctx, ev)
}
