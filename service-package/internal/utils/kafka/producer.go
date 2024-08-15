package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"service-package/internal/domain/dto"
	"time"
)

func KafkaSend(orderReq dto.Response, toTopic string) {
	brokers := []string{"localhost:29092"}
	topic := toTopic

	// Create a new Sarama configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 5                    // Retry up to 5 times to produce the message
	config.Producer.Return.Successes = true

	// Create a new synchronous producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	defer producer.Close()

	// Sending messages
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: 0,
		Key:       sarama.StringEncoder(fmt.Sprintf("%d")),
		Value: sarama.StringEncoder(fmt.Sprintf("OrderType: %s, OrderService: %s, TransactionId: %s, UserId: %s, PackageId: %s, RespCode: %d, RespStatus: %s, RespMessage: %s",
			orderReq.OrderType, orderReq.OrderService, orderReq.TransactionId, orderReq.UserId, orderReq.PackageName, orderReq.RespCode, orderReq.RespStatus, orderReq.RespMessage)),
	}

	log.Printf("Key = %s, Value = %s", msg.Key, msg.Value)
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to partition %d: %v", partition, err)
	} else {
		log.Printf("Message sent to partition %d at offset %d", partition, offset)
	}

	time.Sleep(500 * time.Millisecond) // To simulate delay

}
