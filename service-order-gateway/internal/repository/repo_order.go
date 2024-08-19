package repository

import (
	"context"
	"database/sql"
	"fmt"
	"service-order-gateway/internal/domain"
)

type RepositoryOrder interface {
	CreateOrder(request domain.Order, ctx context.Context) (domain.Order, error)
	UpdateOrder(request domain.Order, ctx context.Context, tx *sql.Tx) (domain.Order, error)
	FindOrderById(transactionId string, ctx context.Context, tx *sql.Tx) (domain.Order, error)
}

type repositoryOrder struct {
	db *sql.DB
}

func (repo *repositoryOrder) FindOrderById(transactionId string, ctx context.Context, tx *sql.Tx) (domain.Order, error) {
	var order domain.Order

	err := tx.QueryRowContext(ctx, "SELECT "+
		"order_type, "+
		"transaction_id, "+
		"user_id, "+
		"product_id, "+
		"quantity, "+
		"total_price, "+
		"status_order, "+
		"created_at, "+
		"updated_at "+
		"FROM "+
		"transaction_orders "+
		"WHERE "+
		"transaction_id = $1", transactionId).Scan(
		&order.OrderType,
		&order.TransactionID,
		&order.UserId,
		&order.ProductId,
		&order.Quantity,
		&order.TotalPrice,
		&order.StatusOrder,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return order, fmt.Errorf("repositoryOrder.FindOrderById: %w", err)
	}

	return order, nil
}

func (repo *repositoryOrder) UpdateOrder(request domain.Order, ctx context.Context, tx *sql.Tx) (domain.Order, error) {
	query := `
        UPDATE transaction_orders
        SET
            order_type = COALESCE(NULLIF($2, ''), order_type),
            user_id = COALESCE(NULLIF($3, ''), user_id),
            product_id = COALESCE(NULLIF($4, ''), product_id),
            quantity = COALESCE(NULLIF($5, 0), quantity),
            total_price = COALESCE(NULLIF($6, 0), total_price),
            status_order = COALESCE(NULLIF($7, ''), status_order),
            updated_at = $8
        WHERE transaction_id = $1
        RETURNING order_type, transaction_id, user_id, product_id, quantity, total_price, status_order, created_at, updated_at;
    `

	var updatedOrder domain.Order
	err := tx.QueryRowContext(ctx, query, request.TransactionID, request.OrderType, request.UserId, request.ProductId, request.Quantity, request.TotalPrice, request.StatusOrder, request.UpdatedAt).Scan(
		&updatedOrder.OrderType,
		&updatedOrder.TransactionID,
		&updatedOrder.UserId,
		&updatedOrder.ProductId,
		&updatedOrder.Quantity,
		&updatedOrder.TotalPrice,
		&updatedOrder.StatusOrder,
		&updatedOrder.CreatedAt,
		&updatedOrder.UpdatedAt,
	)
	if err != nil {
		return domain.Order{}, fmt.Errorf("error updating transaction order: %v", err)
	}

	return updatedOrder, nil
}

func (repo *repositoryOrder) CreateOrder(request domain.Order, ctx context.Context) (domain.Order, error) {
	var order domain.Order

	err := repo.db.QueryRowContext(ctx, "INSERT INTO transaction_orders ("+
		"transaction_id, "+
		"order_type, "+
		"user_id, "+
		"product_id, "+
		"quantity, "+
		"total_price, "+
		"status_order, "+
		"created_at, "+
		"updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) "+
		"RETURNING "+
		"transaction_id, "+
		"order_type, "+
		"user_id, "+
		"product_id, "+
		"quantity, "+
		"total_price, "+
		"status_order, "+
		"created_at, "+
		"updated_at",
		request.TransactionID,
		request.OrderType,
		request.UserId,
		request.ProductId,
		request.Quantity,
		request.TotalPrice,
		request.StatusOrder,
		request.CreatedAt,
		request.UpdatedAt).Scan(
		&order.OrderType,
		&order.TransactionID,
		&order.UserId,
		&order.ProductId,
		&order.Quantity,
		&order.TotalPrice,
		&order.StatusOrder,
		&order.CreatedAt,
		&order.UpdatedAt)
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
