package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"net/http"
	"service-package/internal/domain/dto"
	"service-package/internal/usecase"
	"sync"
)

import (
	"context"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewKafkaConsumer(cfg *Config) (*KafkaConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	version, err := sarama.ParseKafkaVersion(cfg.KafkaVersion)
	if err != nil {
		return nil, err
	}
	saramaConfig.Version = version

	consumer, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroupID, saramaConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaConsumer{
		consumer: consumer,
		topics:   []string{cfg.KafkaTopic},
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (kc *KafkaConsumer) Consume() {
	kc.wg.Add(1)
	go func() {
		defer kc.wg.Done()
		for {
			if err := kc.consumer.Consume(kc.ctx, kc.topics, &ConsumerGroupHandler{}); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if kc.ctx.Err() != nil {
				return
			}
		}
	}()
}

func (kc *KafkaConsumer) Stop() {
	kc.cancel()
	kc.wg.Wait()
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var responseMessageKafka dto.Response
		err := json.Unmarshal(message.Value, &responseMessageKafka)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		fmt.Printf("Received message: %+v\n", responseMessageKafka)

		session.MarkMessage(message, "")

		switch responseMessageKafka.OrderService {
		case "validate_user":
			repoUser := repository.NewRepositoryUser()
			ucUser := usecase.NewUsecaseUser(repoUser)

			sendMessageToValidateUser := dto.Response{
				OrderType:     responseMessageKafka.OrderType,
				OrderService:  "validate_package",
				TransactionId: responseMessageKafka.TransactionId,
				UserId:        responseMessageKafka.UserId,
				PackageName:   responseMessageKafka.PackageName,
				RespCode:      200,
				RespStatus:    "validate user success",
				RespMessage:   "validate user success, next step validate package name",
			}

			_, err := ucUser.ValidateUser(responseMessageKafka.UserId)
			if err != nil {
				log.Printf("Error validating user: %v", err)

				sendMessageToValidateUser.RespCode = http.StatusBadGateway
				sendMessageToValidateUser.RespStatus = http.StatusText(http.StatusBadGateway)
				sendMessageToValidateUser.RespMessage = err.Error()

				KafkaSendMessage(sendMessageToValidateUser, "topic_0")
			}

			KafkaSendMessage(sendMessageToValidateUser, "topic_0")
		}
	}
	return nil
}
