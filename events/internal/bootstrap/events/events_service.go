package events

import (
	"context"
	"encoding/json"

	eventmodel "github.com/sashapremium/events/events/internal/models"
	"github.com/sashapremium/events/events/internal/services/eventsService"
	kproducer "github.com/sashapremium/events/events/internal/services/processors/content_events_processor"
)

type eventProducerAdapter struct {
	p kproducer.Producer
}

func (a eventProducerAdapter) Close() error {
	//TODO implement me
	panic("implement me")
}

func (a eventProducerAdapter) PublishEvent(ctx context.Context, ev *eventmodel.ContentEvent) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return a.p.Send(ctx, []byte(ev.ContentID), b)
}

func InitEventsService(storage eventsService.Storage, producer kproducer.Producer) *eventsService.Service {
	return eventsService.NewService(storage, eventProducerAdapter{p: producer})
}
