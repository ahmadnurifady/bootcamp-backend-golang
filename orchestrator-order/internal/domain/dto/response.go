package dto

type BaseResponse struct {
	StatusCode int    `json:"statuscode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

type MessageKafka struct {
	OrderType      string `json:"orderType"`
	FromService    string `json:"fromService"`
	TakenByService string `json:"takenByService"`
	TransactionId  string `json:"transactionId"`
	UserId         string `json:"userId"`
	ProductId      string `json:"productId"`
	Payload        any    `json:"payload"`
	RespStatus     string `json:"respStatus"`
	RespMessage    string `json:"respMessage"`
	RespCode       int    `json:"respCode"`
}
