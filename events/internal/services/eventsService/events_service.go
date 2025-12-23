package eventsService

import (
	"context"

	eventmodel "github.com/sashapremium/events/events/internal/models"
)

type Storage interface {
	InsertEvents(ctx context.Context, events []*eventmodel.ContentEvent) error
	GetPost(ctx context.Context, id uint64) (*eventmodel.PostInfo, error)
	GetPostAuthorID(ctx context.Context, postID uint64) (uint64, error)
}

type EventProducer interface {
	PublishEvent(ctx context.Context, event *eventmodel.ContentEvent) error
	Close() error
}

type Service struct {
	storage  Storage
	producer EventProducer
}

func NewService(storage Storage, producer EventProducer) *Service {
	return &Service{
		storage:  storage,
		producer: producer,
	}
}
