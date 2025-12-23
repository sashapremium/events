package eventsService

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	eventmodel "github.com/sashapremium/events/events/internal/models"
)

func (s *Service) newEventWithAuthor(
	ctx context.Context,
	typ eventmodel.EventType,
	postID uint64,
	userHash string,
	comment string,
) (*eventmodel.ContentEvent, error) {
	if s.storage == nil {
		return nil, errors.New("storage is nil")
	}

	ev := newEvent(typ, postID, userHash, comment)

	authorID, err := s.storage.GetPostAuthorID(ctx, postID)
	if err != nil {
		return nil, err
	}
	ev.AuthorID = authorID

	return ev, nil
}

func newEvent(typ eventmodel.EventType, postID uint64, userHash, comment string) *eventmodel.ContentEvent {
	return &eventmodel.ContentEvent{
		ContentID: strconv.FormatUint(postID, 10),
		UserHash:  userHash,
		Type:      typ,
		Comment:   comment,
		At:        time.Now().UTC(),
	}
}

func (s *Service) persistAndPublish(ctx context.Context, ev *eventmodel.ContentEvent) error {
	if s.storage == nil {
		return fmt.Errorf("storage is nil")
	}
	if s.producer == nil {
		return fmt.Errorf("producer is nil")
	}

	if err := s.storage.InsertEvents(ctx, []*eventmodel.ContentEvent{ev}); err != nil {
		return err
	}
	if err := s.producer.PublishEvent(ctx, ev); err != nil {
		return err
	}

	return nil
}
