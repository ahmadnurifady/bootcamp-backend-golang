package main

import (
	"fmt"
	"orchestrator-order/internal/delivery/server"
)

func main() {

	err := server.NewServer().Run()
	if err != nil {
		fmt.Println(err)
	}

	//cfg, err := kafka_util.Load()
	//if err != nil {
	//	log.Fatalf("Failed to load configuration: %v", err)
	//}
	//
	//kafkaConsumer, err := kafka_util.NewKafkaConsumer(cfg)
	//if err != nil {
	//	log.Fatalf("Failed to create Kafka consumer: %v", err)
	//}
	//defer kafkaConsumer.Close()
	//
	//go kafkaConsumer.Consume()
	//
	//sigterm := make(chan os.Signal, 1)
	//signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	//<-sigterm
	//
	//kafkaConsumer.Stop()
	//log.Println("Received termination signal, initiating shutdown...")
	//
	//server.NewServer().Run()
	//kafka_util.KafkaSend()
}
