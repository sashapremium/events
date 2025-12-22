package analytics_service_api

import (
	"context"

	"github.com/sashapremium/events/analytics/internal/pb/analytics_api"
)

func (a *AnalyticsServiceAPI) GetAuthorStats(
	ctx context.Context,
	req *analytics_api.GetAuthorStatsRequest,
) (*analytics_api.GetAuthorStatsResponse, error) {
	stats, err := a.svc.GetAuthorStats(ctx, req.AuthorId)
	if err != nil {
		return nil, err
	}

	return &analytics_api.GetAuthorStatsResponse{
		Stats: stats,
	}, nil
}
