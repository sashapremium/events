package analytics_service_api

import (
	"context"

	"github.com/sashapremium/events/analytics/internal/pb/analytics_api"
)

func (a *AnalyticsServiceAPI) GetTop(
	ctx context.Context,
	req *analytics_api.GetTopRequest,
) (*analytics_api.GetTopResponse, error) {
	top, err := a.svc.GetTop(ctx, req.Metric, req.Limit)
	if err != nil {
		return nil, err
	}

	return &analytics_api.GetTopResponse{
		Top: top,
	}, nil
}
