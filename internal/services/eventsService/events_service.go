package eventsService

import (
	"context"

	eventmodel "github.com/sashapremium/events/internal/models"
)

type Storage interface {
	InsertEvents(ctx context.Context, events []*eventmodel.ContentEvent) error
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
