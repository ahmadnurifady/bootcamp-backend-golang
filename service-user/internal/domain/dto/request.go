package dto

type Message struct {
	OrderType     string `json:"orderType"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	PackageId     string `json:"packageId"`
}
