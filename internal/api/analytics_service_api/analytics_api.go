package analytics_service_api

import (
	"context"

	"github.com/sashapremium/events/internal/pb/analytics_api"
	proto_models "github.com/sashapremium/events/internal/pb/models"
)

type AnalyticsService interface {
	GetPostStats(ctx context.Context, postID uint64, fresh bool) (*proto_models.PostStatsModel, error)
	GetTop(ctx context.Context, metric string, limit uint32) (*proto_models.TopModel, error)
	GetAuthorStats(ctx context.Context, authorID string) (*proto_models.AuthorStatsModel, error)
}

type AnalyticsServiceAPI struct {
	analytics_api.UnimplementedAnalyticsServiceServer
	svc AnalyticsService
}

func NewAnalyticsServiceAPI(svc AnalyticsService) *AnalyticsServiceAPI {
	return &AnalyticsServiceAPI{svc: svc}
}
