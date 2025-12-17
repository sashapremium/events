package analyticsService_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sashapremium/events/internal/models"
	analyticsmodels "github.com/sashapremium/events/internal/pb/models"
	"github.com/sashapremium/events/internal/services/analyticsService"
	"github.com/sashapremium/events/internal/services/analyticsService/mocks"
	"github.com/stretchr/testify/require"
)

/*
	Helpers
*/

func defaultCacheMock() *mocks.CacheMock {
	return &mocks.CacheMock{
		GetLastSyncedAtFunc: func(ctx context.Context) (string, bool, error) {
			return "", false, nil
		},
	}
}

/*
	ProcessEvent tests
*/

func TestProcessEvent_View_WithUniqueUser(t *testing.T) {
	cache := defaultCacheMock()
	cache.AddUniqueUserFunc = func(ctx context.Context, postID uint64, userHash string) (bool, error) {
		return true, nil
	}
	cache.IncrTotalsFunc = func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
		require.Equal(t, int64(1), delta.Views)
		require.Equal(t, int64(1), delta.UniqueUsers)
		return nil
	}
	cache.MarkDirtyFunc = func(ctx context.Context, postID uint64) error { return nil }
	cache.IncrTopFunc = func(ctx context.Context, metric string, postID uint64, inc int64) error {
		require.Equal(t, "views", metric)
		return nil
	}

	svc := analyticsService.New(&mocks.StorageMock{}, cache)

	err := svc.ProcessEvent(context.Background(), &models.ContentEvent{
		ContentID: "10",
		Type:      models.EventView,
		UserHash:  "u1",
	})

	require.NoError(t, err)
}

func TestProcessEvent_Like(t *testing.T) {
	cache := defaultCacheMock()
	cache.IncrTotalsFunc = func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
		require.Equal(t, int64(1), delta.Likes)
		return nil
	}
	cache.MarkDirtyFunc = func(ctx context.Context, postID uint64) error { return nil }
	cache.IncrTopFunc = func(ctx context.Context, metric string, postID uint64, inc int64) error {
		require.Equal(t, "likes", metric)
		return nil
	}

	svc := analyticsService.New(&mocks.StorageMock{}, cache)

	require.NoError(t, svc.ProcessEvent(context.Background(), &models.ContentEvent{
		ContentID: "1",
		Type:      models.EventLike,
	}))
}

func TestProcessEvent_Comment(t *testing.T) {
	cache := defaultCacheMock()
	cache.IncrTotalsFunc = func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
		require.Equal(t, int64(1), delta.Comments)
		return nil
	}
	cache.MarkDirtyFunc = func(ctx context.Context, postID uint64) error { return nil }
	cache.IncrTopFunc = func(ctx context.Context, metric string, postID uint64, inc int64) error {
		require.Equal(t, "comments", metric)
		return nil
	}

	svc := analyticsService.New(&mocks.StorageMock{}, cache)

	require.NoError(t, svc.ProcessEvent(context.Background(), &models.ContentEvent{
		ContentID: "1",
		Type:      models.EventComment,
	}))
}

func TestProcessEvent_Repost(t *testing.T) {
	cache := defaultCacheMock()
	cache.IncrTotalsFunc = func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
		require.Equal(t, int64(1), delta.Reposts)
		return nil
	}
	cache.MarkDirtyFunc = func(ctx context.Context, postID uint64) error { return nil }
	cache.IncrTopFunc = func(ctx context.Context, metric string, postID uint64, inc int64) error {
		require.Equal(t, "reposts", metric)
		return nil
	}

	svc := analyticsService.New(&mocks.StorageMock{}, cache)

	require.NoError(t, svc.ProcessEvent(context.Background(), &models.ContentEvent{
		ContentID: "1",
		Type:      models.EventRepost,
	}))
}

func TestProcessEvent_UnknownType(t *testing.T) {
	svc := analyticsService.New(&mocks.StorageMock{}, defaultCacheMock())
	require.NoError(t, svc.ProcessEvent(context.Background(), &models.ContentEvent{
		ContentID: "1",
		Type:      "unknown",
	}))
}

func TestProcessEvent_NilEvent(t *testing.T) {
	svc := analyticsService.New(&mocks.StorageMock{}, defaultCacheMock())
	require.NoError(t, svc.ProcessEvent(context.Background(), nil))
}

func TestProcessEvent_InvalidContentID(t *testing.T) {
	svc := analyticsService.New(&mocks.StorageMock{}, defaultCacheMock())

	err := svc.ProcessEvent(context.Background(), &models.ContentEvent{
		ContentID: "abc",
		Type:      models.EventView,
	})

	require.Error(t, err)
}

func TestProcessEvent_IncrTotalsError(t *testing.T) {
	cache := defaultCacheMock()
	cache.IncrTotalsFunc = func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
		return errors.New("fail")
	}

	svc := analyticsService.New(&mocks.StorageMock{}, cache)

	require.Error(t, svc.ProcessEvent(context.Background(), &models.ContentEvent{
		ContentID: "1",
		Type:      models.EventLike,
	}))
}

/*
	GetPostStats tests
*/

func TestGetPostStats_WithFresh(t *testing.T) {
	storage := &mocks.StorageMock{
		GetPostTotalsFunc: func(ctx context.Context, postID uint64) (analyticsService.PostTotals, error) {
			return analyticsService.PostTotals{Views: 10}, nil
		},
	}

	cache := defaultCacheMock()
	cache.GetDeltaFunc = func(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error) {
		return analyticsService.TotalsDelta{Views: 2}, true, nil
	}

	svc := analyticsService.New(storage, cache)

	out, err := svc.GetPostStats(context.Background(), 1, true)
	require.NoError(t, err)
	require.Equal(t, int64(10), out.Totals.TotalViews)
	require.Equal(t, int64(2), out.FreshTail.Views)
}

func TestGetPostStats_NoFresh(t *testing.T) {
	storage := &mocks.StorageMock{
		GetPostTotalsFunc: func(ctx context.Context, postID uint64) (analyticsService.PostTotals, error) {
			return analyticsService.PostTotals{Views: 5}, nil
		},
	}

	svc := analyticsService.New(storage, defaultCacheMock())

	out, err := svc.GetPostStats(context.Background(), 1, false)
	require.NoError(t, err)
	require.Nil(t, out.FreshTail)
}

func TestGetPostStats_StorageError(t *testing.T) {
	storage := &mocks.StorageMock{
		GetPostTotalsFunc: func(ctx context.Context, postID uint64) (analyticsService.PostTotals, error) {
			return analyticsService.PostTotals{}, errors.New("db error")
		},
	}

	svc := analyticsService.New(storage, defaultCacheMock())
	_, err := svc.GetPostStats(context.Background(), 1, false)
	require.Error(t, err)
}

func TestGetPostStats_GetDeltaError(t *testing.T) {
	storage := &mocks.StorageMock{
		GetPostTotalsFunc: func(ctx context.Context, postID uint64) (analyticsService.PostTotals, error) {
			return analyticsService.PostTotals{}, nil
		},
	}

	cache := defaultCacheMock()
	cache.GetDeltaFunc = func(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error) {
		return analyticsService.TotalsDelta{}, false, errors.New("cache error")
	}

	svc := analyticsService.New(storage, cache)
	_, err := svc.GetPostStats(context.Background(), 1, true)
	require.Error(t, err)
}

/*
	Queries
*/

func TestGetTop_OK(t *testing.T) {
	cache := defaultCacheMock()
	cache.GetTopFunc = func(ctx context.Context, metric string, limit uint32) ([]analyticsService.TopItem, error) {
		return []analyticsService.TopItem{{PostID: 1, Value: 100}}, nil
	}

	svc := analyticsService.New(&mocks.StorageMock{}, cache)

	out, err := svc.GetTop(context.Background(), "views", 10)
	require.NoError(t, err)
	require.Len(t, out.Items, 1)
}

func TestGetTop_InvalidMetric(t *testing.T) {
	svc := analyticsService.New(&mocks.StorageMock{}, defaultCacheMock())
	_, err := svc.GetTop(context.Background(), "bad", 10)
	require.ErrorIs(t, err, analyticsService.ErrInvalidMetric)
}

func TestGetTop_CacheError(t *testing.T) {
	cache := defaultCacheMock()
	cache.GetTopFunc = func(ctx context.Context, metric string, limit uint32) ([]analyticsService.TopItem, error) {
		return nil, errors.New("cache fail")
	}

	svc := analyticsService.New(&mocks.StorageMock{}, cache)
	_, err := svc.GetTop(context.Background(), "views", 10)
	require.Error(t, err)
}

/*
	Author stats
*/

func TestGetAuthorStats_OK(t *testing.T) {
	expected := &analyticsmodels.AuthorStatsModel{}

	storage := &mocks.StorageMock{
		GetAuthorStatsFunc: func(ctx context.Context, authorID string) (*analyticsmodels.AuthorStatsModel, error) {
			return expected, nil
		},
	}

	svc := analyticsService.New(storage, defaultCacheMock())

	out, err := svc.GetAuthorStats(context.Background(), "author-1")
	require.NoError(t, err)
	require.Equal(t, expected, out)
}
func TestFlushOnce_OK(t *testing.T) {
	storage := &mocks.StorageMock{
		UpsertPostTotalsFunc: func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
			require.Equal(t, int64(1), delta.Views)
			return nil
		},
	}

	cache := defaultCacheMock()
	cache.GetDirtyBatchFunc = func(ctx context.Context, limit int) ([]uint64, error) {
		return []uint64{1}, nil
	}
	cache.GetDeltaFunc = func(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error) {
		return analyticsService.TotalsDelta{Views: 1}, true, nil
	}
	cache.ResetDeltaFunc = func(ctx context.Context, postID uint64) error {
		return nil
	}
	cache.SetLastSyncedAtFunc = func(ctx context.Context, ts string) error {
		require.NotEmpty(t, ts)
		return nil
	}

	svc := analyticsService.New(storage, cache)
	svc.FlushOnce(context.Background(), 10)
}
func TestFlushOnce_EmptyOrZeroDelta(t *testing.T) {
	storage := &mocks.StorageMock{}

	cache := defaultCacheMock()
	cache.GetDirtyBatchFunc = func(ctx context.Context, limit int) ([]uint64, error) {
		return []uint64{1}, nil
	}
	cache.GetDeltaFunc = func(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error) {
		return analyticsService.TotalsDelta{}, true, nil
	}
	cache.ResetDeltaFunc = func(ctx context.Context, postID uint64) error {
		return nil
	}

	svc := analyticsService.New(storage, cache)
	svc.FlushOnce(context.Background(), 10)
}
func TestFlushOnce_GetDirtyBatchError(t *testing.T) {
	cache := defaultCacheMock()
	cache.GetDirtyBatchFunc = func(ctx context.Context, limit int) ([]uint64, error) {
		return nil, errors.New("cache error")
	}

	svc := analyticsService.New(&mocks.StorageMock{}, cache)
	svc.FlushOnce(context.Background(), 10)
}
func TestFlushOnce_GetDeltaErrorOrNotOk(t *testing.T) {
	storage := &mocks.StorageMock{}

	cache := defaultCacheMock()
	cache.GetDirtyBatchFunc = func(ctx context.Context, limit int) ([]uint64, error) {
		return []uint64{1, 2}, nil
	}
	cache.GetDeltaFunc = func(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error) {
		if postID == 1 {
			return analyticsService.TotalsDelta{}, false, nil
		}
		return analyticsService.TotalsDelta{}, false, errors.New("cache fail")
	}

	svc := analyticsService.New(storage, cache)
	svc.FlushOnce(context.Background(), 10)
}
func TestFlushOnce_UpsertError(t *testing.T) {
	storage := &mocks.StorageMock{
		UpsertPostTotalsFunc: func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
			return errors.New("db fail")
		},
	}

	cache := defaultCacheMock()
	cache.GetDirtyBatchFunc = func(ctx context.Context, limit int) ([]uint64, error) {
		return []uint64{1}, nil
	}
	cache.GetDeltaFunc = func(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error) {
		return analyticsService.TotalsDelta{Views: 1}, true, nil
	}

	svc := analyticsService.New(storage, cache)
	svc.FlushOnce(context.Background(), 10)
}
