package analytics_service_api

import (
	"context"

	"github.com/sashapremium/events/internal/pb/analytics_api"
)

func (a *AnalyticsServiceAPI) GetPostStats(ctx context.Context, req *analytics_api.GetPostStatsRequest) (*analytics_api.GetPostStatsResponse, error) {
	stats, err := a.svc.GetPostStats(ctx, req.PostId, req.Fresh)
	if err != nil {
		return nil, err
	}

	return &analytics_api.GetPostStatsResponse{Stats: stats}, nil
}
