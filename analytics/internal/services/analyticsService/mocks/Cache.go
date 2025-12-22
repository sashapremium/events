package mocks

import (
	"context"

	"github.com/sashapremium/events/analytics/internal/services/analyticsService"
)

type CacheMock struct {
	IncrTotalsFunc func(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error
	GetTotalsFunc  func(ctx context.Context, postID uint64) (analyticsService.PostTotals, bool, error)

	GetDeltaFunc   func(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error)
	ResetDeltaFunc func(ctx context.Context, postID uint64) error

	AddUniqueUserFunc func(ctx context.Context, postID uint64, userHash string) (bool, error)

	IncrTopFunc func(ctx context.Context, metric string, postID uint64, inc int64) error
	GetTopFunc  func(ctx context.Context, metric string, limit uint32) ([]analyticsService.TopItem, error)

	MarkDirtyFunc     func(ctx context.Context, postID uint64) error
	GetDirtyBatchFunc func(ctx context.Context, limit int) ([]uint64, error)

	SetLastSyncedAtFunc func(ctx context.Context, ts string) error
	GetLastSyncedAtFunc func(ctx context.Context) (string, bool, error)
}

func (m *CacheMock) GetTotals(ctx context.Context, postID uint64) (analyticsService.PostTotals, bool, error) {
	if m.GetTotalsFunc == nil {
		return analyticsService.PostTotals{}, false, nil
	}
	return m.GetTotalsFunc(ctx, postID)
}

func (m *CacheMock) GetDelta(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error) {
	if m.GetDeltaFunc == nil {
		return analyticsService.TotalsDelta{}, false, nil
	}
	return m.GetDeltaFunc(ctx, postID)
}

func (m *CacheMock) ResetDelta(ctx context.Context, postID uint64) error {
	if m.ResetDeltaFunc == nil {
		return nil
	}
	return m.ResetDeltaFunc(ctx, postID)
}

func (m *CacheMock) AddUniqueUser(ctx context.Context, postID uint64, userHash string) (bool, error) {
	if m.AddUniqueUserFunc == nil {
		return false, nil
	}
	return m.AddUniqueUserFunc(ctx, postID, userHash)
}

func (m *CacheMock) IncrTop(ctx context.Context, metric string, postID uint64, inc int64) error {
	if m.IncrTopFunc == nil {
		return nil
	}
	return m.IncrTopFunc(ctx, metric, postID, inc)
}

func (m *CacheMock) GetTop(ctx context.Context, metric string, limit uint32) ([]analyticsService.TopItem, error) {
	if m.GetTopFunc == nil {
		return nil, nil
	}
	return m.GetTopFunc(ctx, metric, limit)
}

func (m *CacheMock) MarkDirty(ctx context.Context, postID uint64) error {
	if m.MarkDirtyFunc == nil {
		return nil
	}
	return m.MarkDirtyFunc(ctx, postID)
}

func (m *CacheMock) GetDirtyBatch(ctx context.Context, limit int) ([]uint64, error) {
	if m.GetDirtyBatchFunc == nil {
		return nil, nil
	}
	return m.GetDirtyBatchFunc(ctx, limit)
}

func (m *CacheMock) SetLastSyncedAt(ctx context.Context, ts string) error {
	if m.SetLastSyncedAtFunc == nil {
		return nil
	}
	return m.SetLastSyncedAtFunc(ctx, ts)
}

func (m *CacheMock) GetLastSyncedAt(ctx context.Context) (string, bool, error) {
	if m.GetLastSyncedAtFunc == nil {
		return "", false, nil
	}
	return m.GetLastSyncedAtFunc(ctx)
}
func (m *CacheMock) IncrTotals(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
	if m.IncrTotalsFunc == nil {
		return nil
	}
	return m.IncrTotalsFunc(ctx, postID, delta)
}
