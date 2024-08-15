package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"orchestrator-order/internal/domain"
	"orchestrator-order/internal/domain/dto"
	"orchestrator-order/internal/usecase"
)

type MessageHandler struct {
	producer *KafkaProducer
	ucTx     usecase.TransactionUsecase
}

func NewMessageHandler(producer *KafkaProducer, ucTx usecase.TransactionUsecase) *MessageHandler {
	return &MessageHandler{producer: producer, ucTx: ucTx}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handleMessage(msg); err != nil {
			log.Printf("Error handling message: %v", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h MessageHandler) handleMessage(msg *sarama.ConsumerMessage) error {
	var messageKafka dto.MessageKafka
	if err := json.Unmarshal(msg.Value, &messageKafka); err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	log.Printf("Received message: %s", string(msg.Value))

	nextStepFlow, err := h.ucTx.GetNextFlowStep(messageKafka.OrderType, messageKafka.FromService, messageKafka.TakenByService)
	if err != nil {
		return fmt.Errorf("getting next flow step: %w", err)
	}
	log.Printf("Next step: %s", nextStepFlow)

	tx := createTransaction(messageKafka)

	if _, err := h.ucTx.CreateTransaction(tx); err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	if tx.StatusService == "success" {
		if err := h.sendToNextStep(messageKafka, nextStepFlow); err != nil {
			return fmt.Errorf("sending to next step: %w", err)
		}
	} else {
		log.Println("Stopping message flow due to unsuccessful response code")
	}

	return nil
}

func createTransaction(msg dto.MessageKafka) domain.Transaction {
	status := "failed"
	if msg.RespCode >= 200 && msg.RespCode <= 299 {
		status = "success"
	}

	return domain.Transaction{
		TransactionID:  msg.TransactionId,
		OrderType:      msg.OrderType,
		UserID:         msg.UserId,
		ProductID:      msg.ProductId,
		Payload:        fmt.Sprintf("%v", msg),
		CurrentService: msg.FromService,
		StatusService:  status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (h MessageHandler) sendToNextStep(msg dto.MessageKafka, nextStep string) error {
	sendMessage, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshalling message: %w", err)
	}

	if err := h.producer.SendMessage(uuid.New().String(), sendMessage, nextStep); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
