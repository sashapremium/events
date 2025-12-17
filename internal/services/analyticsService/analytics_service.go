package analyticsService

import (
	"context"
	"errors"
	"time"

	analyticsmodels "github.com/sashapremium/events/internal/pb/models"
)

type Storage interface {
	UpsertPostTotals(ctx context.Context, postID uint64, delta TotalsDelta) error
	GetPostTotals(ctx context.Context, postID uint64) (PostTotals, error)
	GetAuthorStats(ctx context.Context, authorID string) (*analyticsmodels.AuthorStatsModel, error)
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

type Service struct {
	storage Storage
	cache   Cache

	syncInterval time.Duration
	syncBatch    int
}

func New(storage Storage, cache Cache) *Service {
	return &Service{
		storage:      storage,
		cache:        cache,
		syncInterval: 3 * time.Second,
		syncBatch:    200,
	}
}

type PostTotals struct {
	Views, Likes, Comments, Reposts, UniqueUsers int64
}

type TotalsDelta struct {
	Views, Likes, Comments, Reposts, UniqueUsers int64
}

type TopItem struct {
	PostID uint64
	Value  int64
}

var ErrInvalidMetric = errors.New("некорректная метрика (ожидается views|likes|comments|reposts)")
