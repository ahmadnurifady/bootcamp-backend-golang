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

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	var messageKafka dto.MessageKafka

	for msg := range claim.Messages() {

		err := json.Unmarshal(msg.Value, &messageKafka)
		if err != nil {
			fmt.Println("Error unmarshalling message", err.Error())
		}

		fmt.Println(string(msg.Value))

		changeDataMessage := dto.MessageKafka{
			OrderType:      messageKafka.OrderType,
			FromService:    "validate_payment_topic",
			TakenByService: messageKafka.TakenByService,
			TransactionId:  messageKafka.TransactionId,
			UserId:         messageKafka.UserId,
			ProductId:      messageKafka.ProductId,
			RespStatus:     "validate payment success",
			RespMessage:    "success payment product",
			RespCode:       200,
		}

		sendMessage, err := json.Marshal(changeDataMessage)
		if err != nil {
			fmt.Println("Error marshalling message", err.Error())
		}

		err = h.producer.SendMessage(uuid.New().String(), sendMessage, "topic_0")
		if err != nil {
			fmt.Println("Error sending message", err.Error())
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
