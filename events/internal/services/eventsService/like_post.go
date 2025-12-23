package eventsService

import (
	"context"

	eventmodel "github.com/sashapremium/events/events/internal/models"
)

func (s *Service) LikePost(ctx context.Context, id uint64, userHash string) error {
	if err := s.validateUser(userHash); err != nil {
		return err
	}

	ev, err := s.newEventWithAuthor(ctx, eventmodel.EventLike, id, userHash, "")
	if err != nil {
		return err
	}

	return s.persistAndPublish(ctx, ev)
}
