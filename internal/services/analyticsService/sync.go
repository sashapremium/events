package analyticsService

import (
	"context"
	"time"
)

func (s *Service) RunSyncLoop(ctx context.Context) {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.FlushOnce(ctx, s.syncBatch)
		}
	}
}

// Экспортируем unit-логику
func (s *Service) FlushOnce(ctx context.Context, batch int) {
	postIDs, err := s.cache.GetDirtyBatch(ctx, batch)
	if err != nil || len(postIDs) == 0 {
		return
	}

	var applied int64

	for _, postID := range postIDs {
		delta, ok, err := s.cache.GetDelta(ctx, postID)
		if err != nil || !ok {
			continue
		}

		if delta.Views == 0 && delta.Likes == 0 && delta.Comments == 0 && delta.Reposts == 0 && delta.UniqueUsers == 0 {
			_ = s.cache.ResetDelta(ctx, postID)
			continue
		}

		if err := s.storage.UpsertPostTotals(ctx, postID, delta); err != nil {
			continue
		}
		if err := s.cache.ResetDelta(ctx, postID); err != nil {
			continue
		}

		applied++
	}

	if applied > 0 {
		_ = s.cache.SetLastSyncedAt(ctx, time.Now().UTC().Format(time.RFC3339))
	}
}
