package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"orchestrator-order/internal/config"
	"orchestrator-order/internal/delivery/handler"
	"orchestrator-order/internal/provider/db"
	"orchestrator-order/internal/repository"
	"orchestrator-order/internal/usecase"
	"orchestrator-order/internal/utils/kafka"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	repo     repository.TransactionDetailRepository
	uc       usecase.TransactionUsecase
	engine   *gin.Engine
	host     string
	msg      *kafka.MessageHandler
	cfg      *config.Config
	producer *kafka.KafkaProducer
	consumer *kafka.KafkaConsumer
}

func (s *Server) setupControllers() {

	group := s.engine.Group("/api/v1")
	handler.NewTransactionHandler(s.uc, group).Route()
	//group.Use(middleware.LogMiddleware())
	//authMiddleware := middleware.NewAuthMiddleware(s.jwt)
	//handler.NewHandlerOrder(s.uc, group, authMiddleware).Route()
	//handler.NewHandlerAuth(s.jwt, group).RouteAuth()
}

func (s *Server) setupKafka() error {
	brokers := strings.Split(s.cfg.KafkaBrokers, ",")

	var err error
	s.producer, err = kafka.NewKafkaProducer(brokers, s.cfg.OrchestraTopic)
	if err != nil {
		return fmt.Errorf("error creating kafka producer: %w", err)
	}

	s.msg = kafka.NewMessageHandler(s.producer, s.repo)

	s.consumer, err = kafka.NewKafkaConsumer(brokers, s.cfg.KafkaGroupId, []string{s.cfg.OrchestraTopic}, s.msg)
	if err != nil {
		return fmt.Errorf("error creating Kafka consumer: %w", err)
	}

	return nil
}
func (s *Server) Run() error {
	s.setupControllers()
	if err := s.setupKafka(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := s.consumer.Consume(ctx); err != nil {
			log.Printf("Error consuming Kafka messages: %v", err)
			cancel()
		}
	}()

	srv := &http.Server{
		Addr:    s.host,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server can't run: %v", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	if err := s.consumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}

	if err := s.producer.Close(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	}

	log.Println("Server exiting")
	return nil
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %w", err)
	}

	database, err := db.NewConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %w", err)
	}

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	_, err = kafka.NewKafkaProducer(brokers, cfg.OrchestraTopic)
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}

	producer, err := kafka.NewKafkaProducer(brokers, cfg.OrchestraTopic)
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}

	repo := repository.NewTransactionDetailRepository(database.Conn())
	uc := usecase.NewOrderUsecase(repo, database.Conn(), producer)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		uc:       uc,
		engine:   engine,
		host:     host,
		cfg:      cfg,
		producer: producer,
		repo:     repo,
	}
}
