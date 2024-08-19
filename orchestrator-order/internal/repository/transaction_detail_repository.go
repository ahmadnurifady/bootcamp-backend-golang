package repository

import (
	"database/sql"
	"fmt"
	"orchestrator-order/internal/domain"
)

type TransactionDetailRepository interface {
	CreateTransaction(request domain.Transaction) (domain.Transaction, error)
	FindNextStepFlow(orderType string, fromService string, takenService string) (string, error)
	FindTransaction(transactionId string, currentService string, tx *sql.Tx) (domain.Transaction, error)
	FindNextStepFlowRollback(orderType string, fromService string) (string, error)
}

type transactionDetailRepository struct {
	db *sql.DB
}

func (repo *transactionDetailRepository) FindNextStepFlowRollback(orderType string, fromService string) (string, error) {

	var nextStepRollback string

	err := repo.db.QueryRow("SELECT destination_service "+
		"FROM transaction_step_flow_rollback "+
		"WHERE "+
		"order_type = $1 "+
		"AND "+
		"from_service = $2", orderType, fromService).Scan(&nextStepRollback)

	if err != nil {
		return "", fmt.Errorf("error finding next step flow rollback: %w", err)
	}

	return nextStepRollback, nil
}

func (t *transactionDetailRepository) FindTransaction(
	transactionId string,
	currentService string,
	tx *sql.Tx) (domain.Transaction, error) {

	var transaction domain.Transaction

	err := tx.QueryRow("SELECT "+
		"transaction_id, "+
		"order_type, "+
		"user_id, "+
		"product_id, "+
		"payload, "+
		"current_service, "+
		"status_service, "+
		"created_at, "+
		"updated_at "+
		"FROM "+
		"transactions_detail_logs "+
		"WHERE "+
		"transaction_id = $1 "+
		"AND "+
		"status_service = 'failed' "+
		"AND "+
		"current_service = $2", transactionId, currentService).Scan(
		&transaction.TransactionID,
		&transaction.OrderType,
		&transaction.UserID,
		&transaction.ProductID,
		&transaction.Payload,
		&transaction.DestinationService,
		&transaction.StatusService,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("error getting transaction details: %w", err)
	}

	return transaction, nil
}

func (t *transactionDetailRepository) FindNextStepFlow(
	orderType string,
	fromService string,
	takenService string) (string, error) {

	var nextService string

	err := t.db.QueryRow("SELECT "+
		"current_service "+
		"FROM "+
		"transaction_step_flow "+
		"WHERE "+
		"order_type = $1 "+
		"AND "+
		"from_service = $2 "+
		"AND taken_by_service = $3", orderType, fromService, takenService).Scan(&nextService)

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
			"product_id,  payload, current_service, "+
			"status_service, created_at, updated_at"+
			") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) "+
			"RETURNING transaction_id, order_type, user_id, product_id, payload, current_service, status_service, created_at, updated_at",
		request.TransactionID, request.OrderType, request.UserID, request.ProductID, request.Payload, request.DestinationService, request.StatusService, request.CreatedAt, request.UpdatedAt).Scan(
		&transaction.TransactionID,
		&transaction.OrderType,
		&transaction.UserID,
		&transaction.ProductID,
		&transaction.Payload,
		&transaction.DestinationService,
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
