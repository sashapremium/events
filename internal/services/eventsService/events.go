package eventsService

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	eventmodel "github.com/sashapremium/events/internal/models"
	pbmodels "github.com/sashapremium/events/internal/pb/models"
)

func (s *Service) GetPost(ctx context.Context, id uint64) (*eventmodel.PostInfo, error) {
	if s.storage == nil {
		return nil, errors.New("storage is nil")
	}
	return s.storage.GetPost(ctx, id)
}

func (s *Service) ViewPost(ctx context.Context, id uint64, userHash string) error {
	if err := s.validateUser(userHash); err != nil {
		return err
	}
	return s.persistAndPublish(ctx, newEvent(eventmodel.EventView, id, userHash, ""))
}

func (s *Service) LikePost(ctx context.Context, id uint64, userHash string) error {
	if err := s.validateUser(userHash); err != nil {
		return err
	}
	return s.persistAndPublish(ctx, newEvent(eventmodel.EventLike, id, userHash, ""))
}

func (s *Service) RepostPost(ctx context.Context, id uint64, userHash string) error {
	if err := s.validateUser(userHash); err != nil {
		return err
	}
	return s.persistAndPublish(ctx, newEvent(eventmodel.EventRepost, id, userHash, ""))
}

func (s *Service) AddComment(ctx context.Context, id uint64, userHash, text string) (*pbmodels.CommentModel, error) {
	if err := s.validateComment(userHash, text); err != nil {
		return nil, err
	}

	ev := newEvent(eventmodel.EventComment, id, userHash, text)
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
