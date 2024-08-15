package usecase

import (
	"fmt"
	"service-payment/internal/domain"
	"service-payment/internal/utils/outbond"
)

type UsecasePayment interface {
	ValidatePayment(productId string) (domain.Balance, error)
}

type usecasePayment struct {
}

func (uc usecasePayment) ValidatePayment(productId string) (domain.Balance, error) {
	response, err := outbond.GetPaymentValidation(productId)
	if err != nil {
		fmt.Printf("Error consuming API: %v\n", err)
		return domain.Balance{}, err
	}

	fmt.Printf("Payment validation response: %+v\n", response)

	return response.Data, nil
}

func NewUsecasePayment() UsecasePayment {
	return &usecasePayment{}
}
