package analyticsService

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncSuite struct {
	suite.Suite

	ctx   context.Context
	svc   *Service
	store *MockStorage
	cache *MockCache
}

func TestSyncSuite(t *testing.T) {
	suite.Run(t, new(SyncSuite))
}

func (s *SyncSuite) SetupTest() {
	s.ctx = context.Background()
	s.store = NewMockStorage(s.T())
	s.cache = NewMockCache(s.T())
	s.svc = New(s.store, s.cache)
}

func (s *SyncSuite) TestFlushOnce_GetDirtyBatchError_Returns() {
	s.cache.EXPECT().
		GetDirtyBatch(s.ctx, 10).
		Return(nil, errors.New("redis"))

	s.svc.flushOnce(s.ctx, 10)
}

func (s *SyncSuite) TestFlushOnce_EmptyBatch_Returns() {
	s.cache.EXPECT().
		GetDirtyBatch(s.ctx, 10).
		Return([]uint64{}, nil)

	s.svc.flushOnce(s.ctx, 10)
}

func (s *SyncSuite) TestFlushOnce_GetDeltaNotOk_Skips() {
	postID := uint64(1)

	s.cache.EXPECT().
		GetDirtyBatch(s.ctx, 10).
		Return([]uint64{postID}, nil)

	s.cache.EXPECT().
		GetDelta(s.ctx, postID).
		Return(TotalsDelta{}, false, nil)

	s.svc.flushOnce(s.ctx, 10)
}

func (s *SyncSuite) TestFlushOnce_ZeroDelta_ResetsOnly() {
	postID := uint64(2)

	s.cache.EXPECT().
		GetDirtyBatch(s.ctx, 10).
		Return([]uint64{postID}, nil)

	s.cache.EXPECT().
		GetDelta(s.ctx, postID).
		Return(TotalsDelta{}, true, nil)

	s.cache.EXPECT().
		ResetDelta(s.ctx, postID).
		Return(nil)

	s.svc.flushOnce(s.ctx, 10)
}

func (s *SyncSuite) TestFlushOnce_UpsertError_SkipsResetAndNoLastSynced() {
	postID := uint64(3)
	delta := TotalsDelta{Views: 1}

	s.cache.EXPECT().
		GetDirtyBatch(s.ctx, 10).
		Return([]uint64{postID}, nil)

	s.cache.EXPECT().
		GetDelta(s.ctx, postID).
		Return(delta, true, nil)

	s.store.EXPECT().
		UpsertPostTotals(s.ctx, postID, delta).
		Return(errors.New("pg down"))

	s.svc.flushOnce(s.ctx, 10)
}

func (s *SyncSuite) TestFlushOnce_ResetDeltaError_NoLastSynced() {
	postID := uint64(4)
	delta := TotalsDelta{Likes: 1}

	s.cache.EXPECT().
		GetDirtyBatch(s.ctx, 10).
		Return([]uint64{postID}, nil)

	s.cache.EXPECT().
		GetDelta(s.ctx, postID).
		Return(delta, true, nil)

	s.store.EXPECT().
		UpsertPostTotals(s.ctx, postID, delta).
		Return(nil)

	s.cache.EXPECT().
		ResetDelta(s.ctx, postID).
		Return(errors.New("redis reset failed"))

	s.svc.flushOnce(s.ctx, 10)
}

func (s *SyncSuite) TestFlushOnce_Success_SetsLastSyncedAt() {
	postID := uint64(5)
	delta := TotalsDelta{Comments: 2}

	s.cache.EXPECT().
		GetDirtyBatch(s.ctx, 10).
		Return([]uint64{postID}, nil)

	s.cache.EXPECT().
		GetDelta(s.ctx, postID).
		Return(delta, true, nil)

	s.store.EXPECT().
		UpsertPostTotals(s.ctx, postID, delta).
		Return(nil)

	s.cache.EXPECT().
		ResetDelta(s.ctx, postID).
		Return(nil)

	s.cache.EXPECT().
		SetLastSyncedAt(s.ctx, mock.Anything).
		Return(nil)

	s.svc.flushOnce(s.ctx, 10)
}

func (s *SyncSuite) TestRunSyncLoop_CtxDone_ReturnsImmediately() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// чтобы тест точно не ждал тикер
	s.svc.syncInterval = 24 * time.Hour

	s.svc.RunSyncLoop(ctx)
	require.True(s.T(), true)
}
