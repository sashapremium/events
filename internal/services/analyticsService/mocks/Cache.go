package mocks

import (
	"context"

	"github.com/sashapremium/events/internal/services/analyticsService"
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

func (m *CacheMock) IncrTotals(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
	return m.IncrTotalsFunc(ctx, postID, delta)
}

func (m *CacheMock) GetTotals(ctx context.Context, postID uint64) (analyticsService.PostTotals, bool, error) {
	return m.GetTotalsFunc(ctx, postID)
}

func (m *CacheMock) GetDelta(ctx context.Context, postID uint64) (analyticsService.TotalsDelta, bool, error) {
	return m.GetDeltaFunc(ctx, postID)
}

func (m *CacheMock) ResetDelta(ctx context.Context, postID uint64) error {
	return m.ResetDeltaFunc(ctx, postID)
}

func (m *CacheMock) AddUniqueUser(ctx context.Context, postID uint64, userHash string) (bool, error) {
	return m.AddUniqueUserFunc(ctx, postID, userHash)
}

func (m *CacheMock) IncrTop(ctx context.Context, metric string, postID uint64, inc int64) error {
	return m.IncrTopFunc(ctx, metric, postID, inc)
}

func (m *CacheMock) GetTop(ctx context.Context, metric string, limit uint32) ([]analyticsService.TopItem, error) {
	return m.GetTopFunc(ctx, metric, limit)
}

func (m *CacheMock) MarkDirty(ctx context.Context, postID uint64) error {
	return m.MarkDirtyFunc(ctx, postID)
}

func (m *CacheMock) GetDirtyBatch(ctx context.Context, limit int) ([]uint64, error) {
	return m.GetDirtyBatchFunc(ctx, limit)
}

func (m *CacheMock) SetLastSyncedAt(ctx context.Context, ts string) error {
	return m.SetLastSyncedAtFunc(ctx, ts)
}

func (m *CacheMock) GetLastSyncedAt(ctx context.Context) (string, bool, error) {
	return m.GetLastSyncedAtFunc(ctx)
}
