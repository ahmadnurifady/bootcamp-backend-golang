package dto

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

type BaseResponse struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Data            any    `json:"data"`
}
