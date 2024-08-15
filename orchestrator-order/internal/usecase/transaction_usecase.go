package usecase

import (
	"orchestrator-order/internal/domain"
	"orchestrator-order/internal/repository"
)

type TransactionUsecase interface {
	GetNextFlowStep(orderType string, fromService string, takenService string) (string, error)
	CreateTransaction(request domain.Transaction) (domain.Transaction, error)
}

type transactionUsecase struct {
	repo repository.TransactionDetailRepository
}

func (uc *transactionUsecase) CreateTransaction(request domain.Transaction) (domain.Transaction, error) {

	//sendToRepo := domain.Transaction{
	//	TransactionID:  request.TransactionID,
	//	OrderType:      request.OrderType,
	//	UserId:         request.UserId,
	//	ProductId:      request.ProductId,
	//	Payload:        nil,
	//	CurrentService: "orchestrator",
	//	CreatedAt:      time.Now(),
	//	UpdatedAt:      time.Now(),
	//}

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

func NewOrderUsecase(repo repository.TransactionDetailRepository) TransactionUsecase {

	return &transactionUsecase{repo: repo}
}
