package dto

type RequestLogin struct {
	Username string
	Password string
}

type TransactionRequest struct {
	TransactionID string `json:"transactionId" binding:"required"`
	OrderType     string `json:"orderType" binding:"required"`
	UserId        string `json:"userId" binding:"required"`
	ProductId     string `json:"productId" binding:"required"`
}
