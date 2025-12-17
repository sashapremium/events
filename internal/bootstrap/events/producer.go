package events

import (
	"github.com/sashapremium/events/config"
	"github.com/sashapremium/events/internal/kafka/producer"
)

func InitKafkaProducer(cfg *config.Config) producer.Producer {
	return producer.NewKafkaProducer(cfg.Kafka)
}
