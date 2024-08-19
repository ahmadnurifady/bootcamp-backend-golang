package usecase

import (
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"orchestrator-order/internal/domain"
	"orchestrator-order/internal/domain/dto"
	"orchestrator-order/internal/repository"
	"orchestrator-order/internal/utils/kafka"
	"time"
)

type TransactionUsecase interface {
	GetNextFlowStep(orderType string, fromService string, takenService string) (string, error)
	CreateTransaction(request domain.Transaction) (domain.Transaction, error)
	EditRetryUsecase(request domain.Transaction) (domain.Transaction, error)
	GetNextStepFlowRollbackUsecase(orderType string, fromService string) (string, error)
	//GetPayload(transactionId string, currentService string, statusService string) (domain.Transaction, error)
}

type transactionUsecase struct {
	repo          repository.TransactionDetailRepository
	database      *sql.DB
	kafkaProducer *kafka.KafkaProducer
}

func (uc *transactionUsecase) GetNextStepFlowRollbackUsecase(orderType string, fromService string) (string, error) {
	tx, err := uc.database.Begin()
	if err != nil {
		return "", fmt.Errorf("cannot begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	result, err := uc.repo.FindNextStepFlowRollback(orderType, fromService)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("cannot commit transaction: %v", err)
	}

	return result, nil

}

//func (uc *transactionUsecase) GetPayload(transactionId string, currentService string, statusService string) (domain.Transaction, error) {
//
//	tx, err := uc.database.Begin()
//	if err != nil {
//		return domain.Transaction{}, fmt.Errorf("error starting transaction: %v", err)
//	}
//
//	defer func() {
//		if err != nil {
//			err := tx.Rollback()
//			if err != nil {
//				return
//			}
//		}
//	}()
//
//	findTransaction, err := uc.repo.FindTransaction(transactionId, statusService, currentService, tx)
//	if err != nil {
//		return domain.Transaction{}, err
//	}
//
//	return findTransaction, nil
//}

func (uc *transactionUsecase) EditRetryUsecase(request domain.Transaction) (domain.Transaction, error) {

	tx, err := uc.database.Begin()
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	fmt.Println("request masuk == ", request)

	findTransaction, err := uc.repo.FindTransaction(request.TransactionID, request.DestinationService, tx)
	if err != nil {
		fmt.Println(err)
		return domain.Transaction{}, err
	}

	fmt.Println("findTransaction =", findTransaction)

	if findTransaction.TransactionID == request.TransactionID {
		if request.OrderType != "" {
			findTransaction.OrderType = request.OrderType
		}
		if request.UserID != "" {
			findTransaction.UserID = request.UserID
		}
		if request.ProductID != "" {
			findTransaction.ProductID = request.ProductID
		}

		if request.Payload != "" {
			findTransaction.Payload = request.Payload
		}
		if request.DestinationService != "" {
			findTransaction.DestinationService = request.DestinationService
		}
		findTransaction.UpdatedAt = time.Now()
	}

	messageKafka := dto.MessageKafka{
		OrderType:      findTransaction.OrderType,
		FromService:    "orchestrator",
		TakenByService: "orchestrator",
		TransactionId:  findTransaction.TransactionID,
		UserId:         findTransaction.UserID,
		ProductId:      findTransaction.ProductID,
		Payload:        nil,
		RespStatus:     "edit retry status",
		RespMessage:    "edit retry message",
		RespCode:       0,
	}

	messageBytes, err := json.Marshal(messageKafka)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("error marshalling message: %w", err)
	}

	err = uc.kafkaProducer.SendMessage(uuid.New().String(), messageBytes, findTransaction.DestinationService)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("error sending message: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return domain.Transaction{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return findTransaction, nil
}

func (uc *transactionUsecase) CreateTransaction(request domain.Transaction) (domain.Transaction, error) {

	result, err := uc.repo.CreateTransaction(request)
	if err != nil {
		return domain.Transaction{}, err
	}

	return result, nil
}

func (uc *transactionUsecase) GetNextFlowStep(orderType string, fromService string, takenService string) (string, error) {

	result, err := uc.repo.FindNextStepFlow(orderType, fromService, takenService)
	if err != nil {
		return "", err
	}

	return result, nil
}

func NewOrderUsecase(repo repository.TransactionDetailRepository, database *sql.DB, kafkaProducer *kafka.KafkaProducer) TransactionUsecase {

	return &transactionUsecase{repo: repo, database: database, kafkaProducer: kafkaProducer}
}
