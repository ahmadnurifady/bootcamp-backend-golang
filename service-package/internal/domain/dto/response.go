package dto

type Response struct {
	OrderType     string `json:"orderType"`
	OrderService  string `json:"orderService"`
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	PackageName   string `json:"packageName"`
	RespCode      int    `json:"respCode"`
	RespStatus    string `json:"respStatus"`
	RespMessage   string `json:"respMessage"`
}
