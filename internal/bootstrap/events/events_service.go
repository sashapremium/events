package events

import (
	"context"
	"encoding/json"

	kproducer "github.com/sashapremium/events/internal/kafka/producer"
	eventmodel "github.com/sashapremium/events/internal/models"
	"github.com/sashapremium/events/internal/services/eventsService"
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
