package kafka

import "github.com/IBM/sarama"

type MessageHandler struct {
	producer *KafkaProducer
}

func NewMessageHandler(producer *KafkaProducer) *MessageHandler {
	return &MessageHandler{producer: producer}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		sess.MarkMessage(msg, "")
	}
	return nil
}
