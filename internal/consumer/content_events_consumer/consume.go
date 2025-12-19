package contenteventsconsumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/sashapremium/events/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *ContentEventsConsumer) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           c.groupID,
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("ContentEventsConsumer.ReadMessage error", "error", err.Error())
			continue
		}

		var ev models.ContentEvent
		if err := json.Unmarshal(msg.Value, &ev); err != nil {
			slog.Error("ContentEventsConsumer.Unmarshal error", "error", err.Error())
			continue
		}

		if err := c.processor.Handle(ctx, &ev); err != nil {
			slog.Error("ContentEventsConsumer.Handle error", "error", err.Error())
		}
		slog.Info("ContentEventsConsumer: processed event",
			"type", ev.Type,
			"content_id", ev.ContentID,
			"author_id", ev.AuthorID,
			"user_hash", ev.UserHash,
		)
	}
}
