package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"service-product/internal/domain/dto"
	"service-product/internal/usecase"
)

type MessageHandler struct {
	producer  *Producer
	ucProduct usecase.UsecaseProduct
}

func NewMessageHandler(producer *Producer, ucProduct usecase.UsecaseProduct) *MessageHandler {
	return &MessageHandler{producer: producer, ucProduct: ucProduct}
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

	fmt.Printf("Received message: %s\n", string(msg.Value))

	response, err := h.validateProduct(messageKafka)
	if err != nil {
		return fmt.Errorf("error validating product: %w", err)
	}

	changeDataMessage := h.prepareChangeDataMessage(&messageKafka, response)

	if err := h.sendMessage(changeDataMessage); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	fmt.Printf("Success sending message: %+v\n", changeDataMessage)
	return nil
}

func (h *MessageHandler) validateProduct(messageKafka dto.MessageKafka) (dto.BaseResponse, error) {
	switch messageKafka.OrderType {
	case "beli barang":
		return h.ucProduct.ValidateProduct(messageKafka.ProductId)
	case "beli pulsa":
		return h.ucProduct.ValidateProductCredit(messageKafka.ProductId)
	default:
		return dto.BaseResponse{}, fmt.Errorf("unknown order type: %s", messageKafka.OrderType)
	}
}

func (h *MessageHandler) prepareChangeDataMessage(messageKafka *dto.MessageKafka, response dto.BaseResponse) dto.MessageKafka {
	status := "validate_product_success"
	if response.ResponseCode >= 400 && response.ResponseCode < 500 {
		status = "rollback"
	}

	return dto.MessageKafka{
		OrderType:      messageKafka.OrderType,
		FromService:    "validate_product_topic",
		TakenByService: messageKafka.TakenByService,
		TransactionId:  messageKafka.TransactionId,
		UserId:         messageKafka.UserId,
		ProductId:      messageKafka.ProductId,
		Payload:        response.Data,
		RespStatus:     status,
		RespMessage:    response.ResponseMessage,
		RespCode:       response.ResponseCode,
	}
}

func (h *MessageHandler) sendMessage(message dto.MessageKafka) error {
	sendMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	if err := h.producer.SendMessage(uuid.New().String(), sendMessage, "topic_0"); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
