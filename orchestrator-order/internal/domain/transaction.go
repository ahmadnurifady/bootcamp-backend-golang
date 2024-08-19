package domain

import (
	"time"
)

type Transaction struct {
	TransactionID      string    `json:"transactionId"`
	OrderType          string    `json:"orderType"`
	UserID             string    `json:"userId"`
	ProductID          string    `json:"productId"`
	Payload            string    `json:"payload"`
	DestinationService string    `json:"destinationService"`
	StatusService      string    `json:"statusService"`
	IsFinished         bool      `json:"isFinished"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
