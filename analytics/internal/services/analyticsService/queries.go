package analyticsService

import (
	"context"
	"fmt"
	"strconv"
	"time"

	analyticsmodels "github.com/sashapremium/events/analytics/internal/pb/models"
)

func metricToEventType(metric string) (string, error) {
	switch metric {
	case "views":
		return "view", nil
	case "likes":
		return "like", nil
	case "comments":
		return "comment", nil
	case "reposts":
		return "repost", nil
	default:
		return "", fmt.Errorf("unknown metric: %s", metric)
	}
}

func (s *Service) GetPostStats(ctx context.Context, postID uint64, fresh bool) (*analyticsmodels.PostStatsModel, error) {
	totals, err := s.storage.GetPostTotals(ctx, postID)
	if err != nil {
		return nil, err
	}

	out := &analyticsmodels.PostStatsModel{
		PostId: postID,
		Totals: &analyticsmodels.TotalsModel{
			TotalViews:    int64(totals.Views),
			TotalLikes:    int64(totals.Likes),
			TotalComments: int64(totals.Comments),
			TotalReposts:  int64(totals.Reposts),
			UniqueUsers:   int64(totals.UniqueUsers),
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
	eventType, err := metricToEventType(metric)
	if err != nil {
		return nil, err
	}

	items, err := s.storage.GetTopPostsByType(ctx, eventType, uint64(limit))
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
			Value:  int64(it.Value),
		})
	}

	return out, nil
}

func (s *Service) GetAuthorStats(ctx context.Context, authorID string) (*analyticsmodels.AuthorStatsModel, error) {
	id, err := strconv.ParseUint(authorID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid authorID %q: %w", authorID, err)
	}
	return s.storage.GetAuthorStats(ctx, id)
}
