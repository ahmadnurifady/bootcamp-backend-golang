package usecase

import (
	"fmt"
	"service-payment/internal/domain/dto"
	"service-payment/internal/utils/outbond"
)

type UsecasePayment interface {
	ValidatePayment(productId string) (dto.BaseResponse, error)
	ValidatePaymentKopay(numberPhone string) (dto.BaseResponse, error)
}

type usecasePayment struct {
}

func (uc usecasePayment) ValidatePaymentKopay(numberPhone string) (dto.BaseResponse, error) {
	response, err := outbond.GetPaymentValidationCredit(numberPhone)
	if err != nil {
		return dto.BaseResponse{}, err
	}

	return *response, nil
}

func (uc usecasePayment) ValidatePayment(productId string) (dto.BaseResponse, error) {
	response, err := outbond.GetPaymentValidation(productId)
	if err != nil {
		fmt.Printf("Error consuming API: %v\n", err)
		return dto.BaseResponse{}, err
	}

	fmt.Printf("Payment validation response: %+v\n", response)

	return *response, nil
}

func NewUsecasePayment() UsecasePayment {
	return &usecasePayment{}
}
