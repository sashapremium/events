package eventsService

import (
	"context"
	"time"

	eventmodel "github.com/sashapremium/events/events/internal/models"
	pbmodels "github.com/sashapremium/events/events/internal/pb/models"
)

func (s *Service) AddComment(ctx context.Context, id uint64, userHash, text string) (*pbmodels.CommentModel, error) {
	if err := s.validateComment(userHash, text); err != nil {
		return nil, err
	}

	ev, err := s.newEventWithAuthor(ctx, eventmodel.EventComment, id, userHash, text)
	if err != nil {
		return nil, err
	}

	if err := s.persistAndPublish(ctx, ev); err != nil {
		return nil, err
	}

	return &pbmodels.CommentModel{
		PostId:    id,
		UserHash:  userHash,
		Text:      text,
		CreatedAt: ev.At.Format(time.RFC3339),
	}, nil
}
