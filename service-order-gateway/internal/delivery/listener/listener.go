package listener

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/goccy/go-json"
	"service-order-gateway/internal/domain"
	"service-order-gateway/internal/domain/dto"
	"service-order-gateway/internal/usecase"
	"service-order-gateway/internal/utils/kafka"
	"time"
)

type MessageHandler struct {
	producer *kafka.KafkaProducer
	ucOrder  usecase.UsecaseOrder
}

func NewMessageHandler(producer *kafka.KafkaProducer, ucOrder usecase.UsecaseOrder) *MessageHandler {
	return &MessageHandler{producer: producer, ucOrder: ucOrder}
}

func (h *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.processMessage(msg); err != nil {
			fmt.Printf("Error processing message: %v\n", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h *MessageHandler) processMessage(msg *sarama.ConsumerMessage) error {

	var messageKafka dto.MessageKafka

	if err := json.Unmarshal(msg.Value, &messageKafka); err != nil {
		return fmt.Errorf("error unmarshalling message: %w", err)
	}

	fmt.Println(string(msg.Value))

	order := updateTransactionOrder(messageKafka)

	_, err := h.ucOrder.UpdateOrderUsecase(order, context.Background())
	if err != nil {
		return err
	}

	return nil
}

func updateTransactionOrder(msg dto.MessageKafka) domain.Order {
	status := "FAILED"
	if msg.RespCode >= 200 && msg.RespCode <= 299 {
		status = "COMPLETED"
	} else if msg.RespStatus == "rollback" {
		status = "ROLLBACK SUCCESS"
	}

	order := domain.Order{
		OrderType:     msg.OrderType,
		TransactionID: msg.TransactionId,
		UserId:        msg.UserId,
		ProductId:     msg.ProductId,
		StatusOrder:   status,
		UpdatedAt:     time.Now(),
	}

	return order

}
