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
	"orchestrator-order/internal/repository"
)

type MessageHandler struct {
	producer *KafkaProducer
	repoTx   repository.TransactionDetailRepository
}

func NewMessageHandler(producer *KafkaProducer, repoTx repository.TransactionDetailRepository) *MessageHandler {
	return &MessageHandler{producer: producer, repoTx: repoTx}
}

func (h *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handleMessage(msg); err != nil {
			log.Printf("Error handling message: %v", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h *MessageHandler) handleMessage(msg *sarama.ConsumerMessage) error {
	var messageKafka dto.MessageKafka
	if err := json.Unmarshal(msg.Value, &messageKafka); err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	log.Printf("Received message: %s", string(msg.Value))

	nextStepFlow, err := h.repoTx.FindNextStepFlow(messageKafka.OrderType, messageKafka.FromService, messageKafka.TakenByService)
	if err != nil {
		return fmt.Errorf("getting next flow step: %w", err)
	}
	log.Printf("Next step is == %s", nextStepFlow)

	tx := createTransaction(messageKafka)

	if _, err := h.repoTx.CreateTransaction(tx); err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	if err := h.processNextStep(messageKafka, tx.StatusService, nextStepFlow); err != nil {
		return err
	}

	return nil
}

func createTransaction(msg dto.MessageKafka) domain.Transaction {
	status := "failed"
	if msg.RespCode >= 200 && msg.RespCode <= 299 {
		status = "success"
	}

	messageKafka, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
	}

	return domain.Transaction{
		TransactionID:      msg.TransactionId,
		OrderType:          msg.OrderType,
		UserID:             msg.UserId,
		ProductID:          msg.ProductId,
		Payload:            string(messageKafka),
		DestinationService: msg.FromService,
		StatusService:      status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func (h *MessageHandler) processNextStep(msg dto.MessageKafka, statusService string, nextStepFlow string) error {
	if statusService == "success" && nextStepFlow != "finish_process_order" {
		if err := h.sendToNextStep(msg, nextStepFlow); err != nil {
			return fmt.Errorf("sending to next step: %w", err)
		}
	} else {
		if msg.RespStatus == "rollback" {
			nextStepFlow, err := h.repoTx.FindNextStepFlowRollback(msg.OrderType, msg.FromService)
			if err != nil {
				return fmt.Errorf("getting next flow step: %w", err)
			}
			log.Printf("Next step is == %s", nextStepFlow)

			if err := h.sendToNextStep(msg, nextStepFlow); err != nil {
				return fmt.Errorf("sending to next step rollback: %w", err)
			}
		} else {
			if err := h.sendToNextStep(msg, "order_topic"); err != nil {
				return fmt.Errorf("sending to order_topic cause process failed: %w", err)
			}
			log.Println("Stopping message flow")
		}
	}
	return nil
}

func (h *MessageHandler) sendToNextStep(msg dto.MessageKafka, nextStep string) error {
	sendMessage, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshalling message: %w", err)
	}

	if err := h.producer.SendMessage(uuid.New().String(), sendMessage, nextStep); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
