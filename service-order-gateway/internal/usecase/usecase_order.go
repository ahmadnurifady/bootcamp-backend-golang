package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"service-order-gateway/internal/domain"
	"service-order-gateway/internal/domain/dto"
	"service-order-gateway/internal/repository"
	"service-order-gateway/internal/utils/kafka"
	"time"
)

type UsecaseOrder interface {
	CreateOrderUsecase(request dto.OrderRequest, ctx context.Context) (domain.Order, error)
}

type usecaseOrder struct {
	repo              repository.RepositoryOrder
	orchestraProducer *kafka.KafkaProducer
}

func (uc *usecaseOrder) CreateOrderUsecase(request dto.OrderRequest, ctx context.Context) (domain.Order, error) {

	order := domain.Order{
		OrderType:     request.OrderType,
		TransactionID: request.TransactionID,
		UserId:        request.UserId,
		ProductId:     request.ProductId,
		Quantity:      request.Quantity,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	result, err := uc.repo.CreateOrder(order, ctx)
	if err != nil {
		return domain.Order{}, err
	}

	sendMessage := dto.MessageKafka{
		OrderType:      request.OrderType,
		FromService:    "start_order_product",
		TakenByService: "orchestrator",
		TransactionId:  request.TransactionID,
		UserId:         request.UserId,
		ProductId:      request.ProductId,
		Payload:        "",
		RespStatus:     "start process order",
		RespMessage:    "",
		RespCode:       200,
	}

	sendMessageByte, err := json.Marshal(sendMessage)
	if err != nil {
		fmt.Println("error parsing send message to orchestrator topic", err)
		return domain.Order{}, err
	}

	err = uc.orchestraProducer.SendMessage(uuid.New().String(), sendMessageByte)
	if err != nil {
		fmt.Println("error sending message to orchestrator topic", err)
		return domain.Order{}, err
	}

	return result, nil
}

func NewUsecaseOrder(repo repository.RepositoryOrder, orchestraProducer *kafka.KafkaProducer) UsecaseOrder {
	return &usecaseOrder{
		repo:              repo,
		orchestraProducer: orchestraProducer,
	}
}
