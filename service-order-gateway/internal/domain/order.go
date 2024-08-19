package domain

import "time"

type Order struct {
	OrderType     string    `json:"orderType" binding:"required" valo:"notblank,sizeMin=2,sizeMax=100"`
	TransactionID string    `json:"transactionId" binding:"required" valo:"notblank,sizeMin=2,sizeMax=100"`
	UserId        string    `json:"userId" binding:"required" valo:"notblank,sizeMin=2,sizeMax=100"`
	ProductId     string    `json:"productId" binding:"required" valo:"notblank,sizeMin=2,sizeMax=100"`
	Quantity      int       `json:"quantity" binding:"required"`
	TotalPrice    int       `json:"totalPrice" binding:"required"`
	StatusOrder   string    `json:"statusOrder" binding:"required"`
	CreatedAt     time.Time `json:"createdAt" binding:"required" valo:"notblank,sizeMin=2,sizeMax=200"`
	UpdatedAt     time.Time `json:"updatedAt" binding:"required" valo:"notblank,sizeMin=2,sizeMax=200"`
}
