package usecase

import (
	"context"
	"database/sql"
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
	UpdateOrderUsecase(request domain.Order, ctx context.Context) (domain.Order, error)
}

type usecaseOrder struct {
	repo              repository.RepositoryOrder
	orchestraProducer *kafka.KafkaProducer
	database          *sql.DB
}

func (uc *usecaseOrder) UpdateOrderUsecase(request domain.Order, ctx context.Context) (domain.Order, error) {

	tx, err := uc.database.BeginTx(ctx, nil)
	if err != nil {
		return domain.Order{}, fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	findOrder, err := uc.repo.FindOrderById(request.TransactionID, ctx, tx)
	if err != nil {
		return domain.Order{}, err
	}

	fmt.Println("hasil find order ==", findOrder)

	if findOrder.TransactionID == request.TransactionID && findOrder.OrderType == request.OrderType {
		if request.OrderType != "" {
			findOrder.OrderType = request.OrderType
		}
		if request.TransactionID != "" {
			findOrder.TransactionID = request.TransactionID
		}
		if request.UserId != "" {
			findOrder.UserId = request.UserId
		}
		if request.ProductId != "" {
			findOrder.ProductId = request.ProductId
		}
		if request.Quantity != 0 {
			findOrder.Quantity = request.Quantity
		}
		if request.TotalPrice != 0 {
			findOrder.TotalPrice = request.TotalPrice
		}
		if request.StatusOrder != "" {
			findOrder.StatusOrder = request.StatusOrder
		}
		findOrder.UpdatedAt = time.Now()
	}

	updateOrder, err := uc.repo.UpdateOrder(findOrder, ctx, tx)
	if err != nil {
		return domain.Order{}, err
	}

	fmt.Println("update order ==", updateOrder)

	err = tx.Commit()
	if err != nil {
		return domain.Order{}, fmt.Errorf("error committing transaction %v", err)
	}

	return updateOrder, nil

}

func (uc *usecaseOrder) CreateOrderUsecase(request dto.OrderRequest, ctx context.Context) (domain.Order, error) {

	order := domain.Order{
		OrderType:     request.OrderType,
		TransactionID: request.TransactionID,
		UserId:        request.UserId,
		ProductId:     request.ProductId,
		Quantity:      request.Quantity,
		TotalPrice:    request.Quantity * 100000,
		StatusOrder:   "SUCCESS",
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
		Payload:        fmt.Sprintf("%v", request),
		RespStatus:     "start process order",
		RespMessage:    "start process order",
		RespCode:       200,
	}

	if request.OrderType == "beli pulsa" {
		sendMessage.FromService = "start_order_pulsa"
	}

	fmt.Println("hasil message kafka == ", sendMessage)

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

func NewUsecaseOrder(repo repository.RepositoryOrder, orchestraProducer *kafka.KafkaProducer, database *sql.DB) UsecaseOrder {
	return &usecaseOrder{
		repo:              repo,
		orchestraProducer: orchestraProducer,
		database:          database,
	}
}
