package mocks

import (
	"context"

	"github.com/sashapremium/events/internal/models"
	analyticsmodels "github.com/sashapremium/events/internal/pb/models"
	"github.com/sashapremium/events/internal/services/analyticsService"
)

type StorageMock struct {
	InsertEventsFunc      func(ctx context.Context, events []*models.ContentEvent) error
	GetPostTotalsFunc     func(ctx context.Context, postID uint64) (analyticsService.PostTotals, error)
	GetTopPostsByTypeFunc func(ctx context.Context, eventType string, limit uint64) ([]analyticsService.TopItem, error)
	GetAuthorStatsFunc    func(ctx context.Context, authorID uint64) (*analyticsmodels.AuthorStatsModel, error)
}

func (m *StorageMock) InsertEvents(ctx context.Context, events []*models.ContentEvent) error {
	if m.InsertEventsFunc == nil {
		return nil
	}
	return m.InsertEventsFunc(ctx, events)
}

func (m *StorageMock) GetPostTotals(ctx context.Context, postID uint64) (analyticsService.PostTotals, error) {
	if m.GetPostTotalsFunc == nil {
		return analyticsService.PostTotals{}, nil
	}
	return m.GetPostTotalsFunc(ctx, postID)
}

func (m *StorageMock) GetTopPostsByType(
	ctx context.Context,
	eventType string,
	limit uint64,
) ([]analyticsService.TopItem, error) {
	if m.GetTopPostsByTypeFunc == nil {
		return nil, nil
	}
	return m.GetTopPostsByTypeFunc(ctx, eventType, limit)
}

func (m *StorageMock) GetAuthorStats(
	ctx context.Context,
	authorID uint64,
) (*analyticsmodels.AuthorStatsModel, error) {
	if m.GetAuthorStatsFunc == nil {
		return nil, nil
	}
	return m.GetAuthorStatsFunc(ctx, authorID)
}
