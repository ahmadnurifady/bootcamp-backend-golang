package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"service-payment/internal/domain/dto"
	"service-payment/internal/usecase"
)

type MessageHandler struct {
	producer  *Producer
	ucPayment usecase.UsecasePayment
}

func NewMessageHandler(producer *Producer, ucPayment usecase.UsecasePayment) *MessageHandler {
	return &MessageHandler{producer: producer, ucPayment: ucPayment}
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

	var responseOutbond dto.BaseResponse
	var err error

	switch messageKafka.OrderType {
	case "beli barang":
		responseOutbond, err = h.ucPayment.ValidatePayment(messageKafka.UserId)
	case "beli pulsa":
		responseOutbond, err = h.ucPayment.ValidatePaymentKopay(messageKafka.UserId)
	}

	if err != nil {
		return fmt.Errorf("error validating payment: %w", err)
	}

	changeDataMessage := h.prepareChangeDataMessage(&messageKafka, &responseOutbond)

	if err := h.sendMessage(changeDataMessage); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	fmt.Printf("Success sending message: %+v\n", changeDataMessage)
	return nil
}

func (h *MessageHandler) prepareChangeDataMessage(messageKafka *dto.MessageKafka, responseOutbond *dto.BaseResponse) dto.MessageKafka {
	status := "validate_payment_success"
	if responseOutbond.ResponseCode >= 400 && responseOutbond.ResponseCode < 500 {
		status = "rollback"
	}

	return dto.MessageKafka{
		OrderType:      messageKafka.OrderType,
		FromService:    "validate_payment_topic",
		TakenByService: messageKafka.TakenByService,
		TransactionId:  messageKafka.TransactionId,
		UserId:         messageKafka.UserId,
		ProductId:      messageKafka.ProductId,
		Payload:        responseOutbond.Data,
		RespStatus:     status,
		RespMessage:    responseOutbond.ResponseMessage,
		RespCode:       responseOutbond.ResponseCode,
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
