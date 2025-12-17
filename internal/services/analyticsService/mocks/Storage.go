package mocks

import (
	"context"

	analyticsmodels "github.com/sashapremium/events/internal/pb/models"
	"github.com/sashapremium/events/internal/services/analyticsService"
)

type StorageMock struct {
	UpsertPostTotalsFunc func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error
	GetPostTotalsFunc    func(ctx context.Context, postID uint64) (analyticsService.PostTotals, error)
	GetAuthorStatsFunc   func(ctx context.Context, authorID string) (*analyticsmodels.AuthorStatsModel, error)
}

func (m *StorageMock) UpsertPostTotals(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
	return m.UpsertPostTotalsFunc(ctx, postID, delta)
}

func (m *StorageMock) GetPostTotals(ctx context.Context, postID uint64) (analyticsService.PostTotals, error) {
	return m.GetPostTotalsFunc(ctx, postID)
}

func (m *StorageMock) GetAuthorStats(ctx context.Context, authorID string) (*analyticsmodels.AuthorStatsModel, error) {
	return m.GetAuthorStatsFunc(ctx, authorID)
}
