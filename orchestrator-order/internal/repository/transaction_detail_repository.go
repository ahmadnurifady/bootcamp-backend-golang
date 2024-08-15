package repository

import (
	"database/sql"
	"fmt"
	"orchestrator-order/internal/domain"
)

type TransactionDetailRepository interface {
	CreateTransaction(request domain.Transaction) (domain.Transaction, error)
	FindNextStepFlow(orderType string, fromService string, takenService string) (string, error)
}

type transactionDetailRepository struct {
	db *sql.DB
}

func (t *transactionDetailRepository) FindNextStepFlow(orderType string, fromService string, takenService string) (string, error) {
	var nextService string

	err := t.db.QueryRow("SELECT current_service FROM transaction_step_flow WHERE order_type = $1 AND from_service = $2 AND taken_by_service = $3", orderType, fromService, takenService).Scan(&nextService)
	if err != nil {
		return "", err
	}
	return nextService, nil
}

func (t *transactionDetailRepository) CreateTransaction(request domain.Transaction) (domain.Transaction, error) {
	var transaction domain.Transaction

	err := t.db.QueryRow(
		"INSERT INTO transactions_detail_logs ("+
			"transaction_id, order_type, user_id, "+
			"product_id, quantity, total_price,  payload, current_service, "+
			"status_service, created_at, updated_at"+
			") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) "+
			"RETURNING transaction_id, order_type, user_id, product_id, quantity, total_price, payload, current_service, status_service, created_at, updated_at",
		request.TransactionID, request.OrderType, request.UserID, request.ProductID, request.Quantity, request.TotalPrice, request.Payload, request.CurrentService, request.StatusService, request.CreatedAt, request.UpdatedAt).Scan(
		&transaction.TransactionID,
		&transaction.OrderType,
		&transaction.UserID,
		&transaction.ProductID,
		&transaction.Quantity,
		&transaction.TotalPrice,
		&transaction.Payload,
		&transaction.CurrentService,
		&transaction.StatusService,
		&transaction.CreatedAt,
		&transaction.UpdatedAt)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("error creating transaction: %w", err)
	}

	return transaction, nil
}

func NewTransactionDetailRepository(db *sql.DB) TransactionDetailRepository {
	return &transactionDetailRepository{
		db: db,
	}
}
