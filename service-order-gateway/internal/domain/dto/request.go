package dto

import "time"

type OrderRequest struct {
	OrderType     string `json:"orderType" binding:"required" valo:"notblank,sizeMin=2,sizeMax=100"`
	TransactionID string `json:"transactionId" binding:"required" valo:"notblank,sizeMin=2,sizeMax=100"`
	UserId        string `json:"userId" binding:"required" valo:"notblank,sizeMin=2,sizeMax=100"`
	ProductId     string `json:"productId" binding:"required" valo:"notblank,sizeMin=2,sizeMax=100"`
	Quantity      int    `json:"quantity" binding:"required" valo:"min=1,max=100"`
}

type RequestLog struct {
	AcessTime time.Time
	Latency   time.Duration
	ClientIP  string
	Method    string
	Code      int
	Path      string
	UserAgent string
}
