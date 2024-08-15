package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"service-user/internal/domain"
	"service-user/internal/domain/dto"
	"service-user/internal/usecase"
)

type MessageHandler struct {
	producer *Producer
	ucUser   usecase.UsecaseUser
}

func NewMessageHandler(producer *Producer, ucUser usecase.UsecaseUser) *MessageHandler {
	return &MessageHandler{producer: producer, ucUser: ucUser}
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

	responseOutbond, err := h.ucUser.ValidateUser(messageKafka.UserId)
	if err != nil {
		return fmt.Errorf("error validating user: %w", err)
	}

	changeDataMessage := h.prepareChangeDataMessage(&messageKafka, &responseOutbond)

	if err := h.sendMessage(changeDataMessage); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func (h *MessageHandler) prepareChangeDataMessage(messageKafka *dto.MessageKafka, responseOutbond *dto.BaseResponse) dto.MessageKafka {
	status := "validate_user_success"
	payload := responseOutbond.Data
	if responseOutbond.ResponseCode >= 400 && responseOutbond.ResponseCode < 500 {
		status = "validate_user_failed"
		payload = domain.User{}
	}

	return dto.MessageKafka{
		OrderType:      messageKafka.OrderType,
		FromService:    "validate_user_topic",
		TakenByService: messageKafka.TakenByService,
		TransactionId:  messageKafka.TransactionId,
		UserId:         messageKafka.UserId,
		ProductId:      messageKafka.ProductId,
		Payload:        payload,
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
