package analyticsService

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestRunSyncLoop_TickerTriggersFlush(t *testing.T) {
	store := NewMockStorage(t)
	cache := NewMockCache(t)
	svc := New(store, cache)
	svc.syncInterval = 10 * time.Millisecond
	svc.syncBatch = 10

	cache.EXPECT().
		GetDirtyBatch(mock.Anything, 10).
		Return([]uint64{}, nil).
		Maybe()

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	svc.RunSyncLoop(ctx)
}
