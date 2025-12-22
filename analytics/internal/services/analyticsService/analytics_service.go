package analyticsService

import (
	"context"

	"github.com/sashapremium/events/analytics/internal/models"
	analyticsmodels "github.com/sashapremium/events/analytics/internal/pb/models"
)

type PostTotals struct {
	Views       uint64
	Likes       uint64
	Comments    uint64
	Reposts     uint64
	UniqueUsers uint64
}

type TopItem struct {
	PostID uint64
	Value  uint64
}

type Storage interface {
	InsertEvents(ctx context.Context, events []*models.ContentEvent) error
	GetPostTotals(ctx context.Context, postID uint64) (PostTotals, error)
	GetTopPostsByType(ctx context.Context, eventType string, limit uint64) ([]TopItem, error)
	GetAuthorStats(ctx context.Context, authorID uint64) (*analyticsmodels.AuthorStatsModel, error)
}

type Cache interface {
	IncrTotals(ctx context.Context, postID uint64, delta TotalsDelta) error
	GetTotals(ctx context.Context, postID uint64) (PostTotals, bool, error)

	GetDelta(ctx context.Context, postID uint64) (TotalsDelta, bool, error)
	ResetDelta(ctx context.Context, postID uint64) error

	AddUniqueUser(ctx context.Context, postID uint64, userHash string) (bool, error)

	IncrTop(ctx context.Context, metric string, postID uint64, inc int64) error
	GetTop(ctx context.Context, metric string, limit uint32) ([]TopItem, error)

	MarkDirty(ctx context.Context, postID uint64) error
	GetDirtyBatch(ctx context.Context, limit int) ([]uint64, error)

	SetLastSyncedAt(ctx context.Context, ts string) error
	GetLastSyncedAt(ctx context.Context) (string, bool, error)
}
