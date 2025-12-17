package analyticsService

import (
	"context"
	"time"

	analyticsmodels "github.com/sashapremium/events/internal/pb/models"
)

func (s *Service) GetPostStats(ctx context.Context, postID uint64, fresh bool) (*analyticsmodels.PostStatsModel, error) {
	totals, err := s.storage.GetPostTotals(ctx, postID)
	if err != nil {
		return nil, err
	}

	out := &analyticsmodels.PostStatsModel{
		PostId: postID,
		Totals: &analyticsmodels.TotalsModel{
			TotalViews:    totals.Views,
			TotalLikes:    totals.Likes,
			TotalComments: totals.Comments,
			TotalReposts:  totals.Reposts,
			UniqueUsers:   totals.UniqueUsers,
		},
	}

	if ts, ok, err := s.cache.GetLastSyncedAt(ctx); err == nil && ok {
		out.LastSyncedAt = ts
	} else {
		out.LastSyncedAt = time.Time{}.UTC().Format(time.RFC3339)
	}

	if fresh {
		d, ok, err := s.cache.GetDelta(ctx, postID)
		if err != nil {
			return nil, err
		}
		if ok {
			out.FreshTail = &analyticsmodels.FreshTailModel{
				Views:       d.Views,
				Likes:       d.Likes,
				Comments:    d.Comments,
				Reposts:     d.Reposts,
				UniqueUsers: d.UniqueUsers,
			}
		} else {
			out.FreshTail = &analyticsmodels.FreshTailModel{}
		}
	}

	return out, nil
}

func (s *Service) GetTop(ctx context.Context, metric string, limit uint32) (*analyticsmodels.TopModel, error) {
	if err := s.validateMetric(metric); err != nil {
		return nil, err
	}

	items, err := s.cache.GetTop(ctx, metric, limit)
	if err != nil {
		return nil, err
	}

	out := &analyticsmodels.TopModel{
		Metric: metric,
		Items:  make([]*analyticsmodels.TopItemModel, 0, len(items)),
	}

	for _, it := range items {
		out.Items = append(out.Items, &analyticsmodels.TopItemModel{
			PostId: it.PostID,
			Value:  it.Value,
		})
	}

	return out, nil
}

func (s *Service) GetAuthorStats(ctx context.Context, authorID string) (*analyticsmodels.AuthorStatsModel, error) {
	return s.storage.GetAuthorStats(ctx, authorID)
}
