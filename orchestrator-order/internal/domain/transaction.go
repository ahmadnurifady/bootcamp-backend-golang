package domain

import "time"

type Transaction struct {
	TransactionID  string    `json:"transactionId"`
	OrderType      string    `json:"orderType"`
	UserID         string    `json:"userId"`
	ProductID      string    `json:"productId"`
	Quantity       int       `json:"quantity"`
	TotalPrice     int       `json:"totalPrice"`
	Payload        string    `json:"payload"`
	CurrentService string    `json:"currentService"`
	StatusService  string    `json:"statusService"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
