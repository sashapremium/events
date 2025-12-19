package analyticsService

import (
	"context"
	"strconv"

	"github.com/sashapremium/events/internal/models"
)

func (s *Service) ProcessEvent(ctx context.Context, ev *models.ContentEvent) error {
	if ev == nil {
		return nil
	}

	if err := s.storage.InsertEvents(ctx, []*models.ContentEvent{ev}); err != nil {
		return err
	}

	postID, err := strconv.ParseUint(ev.ContentID, 10, 64)
	if err != nil {
		return err
	}

	var delta TotalsDelta
	var metric string

	switch ev.Type {
	case models.EventView:
		delta.Views = 1
		metric = "views"
		if ev.UserHash != "" {
			added, err := s.cache.AddUniqueUser(ctx, postID, ev.UserHash)
			if err != nil {
				return err
			}
			if added {
				delta.UniqueUsers = 1
			}
		}

	case models.EventLike:
		delta.Likes = 1
		metric = "likes"

	case models.EventComment:
		delta.Comments = 1
		metric = "comments"

	case models.EventRepost:
		delta.Reposts = 1
		metric = "reposts"
	}

	if err := s.cache.IncrTotals(ctx, postID, delta); err != nil {
		return err
	}
	if err := s.cache.MarkDirty(ctx, postID); err != nil {
		return err
	}
	if metric != "" {
		if err := s.cache.IncrTop(ctx, metric, postID, 1); err != nil {
			return err
		}
	}

	return nil
}
