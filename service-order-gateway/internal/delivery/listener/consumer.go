package listener

import (
	"context"
	"github.com/IBM/sarama"
	"log"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  sarama.ConsumerGroupHandler
}

func NewKafkaConsumer(
	brokers []string, groupID string,
	topics []string,
	messageHandler *MessageHandler) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	log.Printf("Creating Kafka consumer with brokers: %v, groupID: %s, and topics: %v", brokers, groupID, topics)
	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		topics:   topics,
		handler:  messageHandler,
	}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context) error {
	for {
		err := kc.consumer.Consume(ctx, kc.topics, kc.handler)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
