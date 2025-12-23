package contenteventsprocessor

import (
	"context"
	"fmt"
	"time"

	"github.com/sashapremium/events/events/config"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Send(ctx context.Context, key []byte, value []byte) error
	Close() error
}

type KafkaProducer struct {
	w *kafka.Writer
}

func NewKafkaProducer(cfg config.KafkaConfig) *KafkaProducer {
	return &KafkaProducer{
		w: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			Balancer:     &kafka.Hash{},
			BatchTimeout: 20 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *KafkaProducer) Send(ctx context.Context, key []byte, value []byte) error {
	if err := p.w.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	}); err != nil {
		return fmt.Errorf("kafka write: %w", err)
	}
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.w.Close()
}
