package main

import (
	"log"
	"os"
	"os/signal"
	"service-package/internal/utils/kafka"
	"syscall"
)

func main() {

	cfg, err := kafka.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	kafkaConsumer, err := kafka.NewKafkaConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	go kafkaConsumer.Consume()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	kafkaConsumer.Stop()
	log.Println("Received termination signal, initiating shutdown...")

}
