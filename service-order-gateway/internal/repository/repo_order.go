package repository

import (
	"context"
	"database/sql"
	"service-order-gateway/internal/domain"
)

type RepositoryOrder interface {
	CreateOrder(request domain.Order, ctx context.Context) (domain.Order, error)
}

type repositoryOrder struct {
	db *sql.DB
}

func (repo *repositoryOrder) CreateOrder(request domain.Order, ctx context.Context) (domain.Order, error) {
	var order domain.Order

	err := repo.db.QueryRowContext(ctx, "INSERT INTO transaction_orders (transaction_id, order_type, user_id, product_id, quantity ,created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING transaction_id, order_type, user_id, product_id, quantity ,created_at, updated_at", request.TransactionID, request.OrderType, request.UserId, request.ProductId, request.Quantity, request.CreatedAt, request.UpdatedAt).Scan(&order.OrderType, &order.TransactionID, &order.UserId, &order.ProductId, &order.Quantity, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil

}

func NewRepositoryOrder(db *sql.DB) RepositoryOrder {
	return &repositoryOrder{
		db: db,
	}
}
