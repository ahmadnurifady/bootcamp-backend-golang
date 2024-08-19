package usecase

import (
	"fmt"
	"service-product/internal/domain/dto"
	"service-product/internal/utils/outbond"
)

type UsecaseProduct interface {
	ValidateProduct(productId string) (dto.BaseResponse, error)
	ValidateProductCredit(productId string) (dto.BaseResponse, error)
}

type usecaseProduct struct {
}

func (uc usecaseProduct) ValidateProductCredit(productId string) (dto.BaseResponse, error) {
	response, err := outbond.GetProductCreditValidation(productId)
	if err != nil {
		return dto.BaseResponse{}, err
	}

	fmt.Printf("Product credit validation response: %+v\n", response)

	return *response, nil
}

func (uc usecaseProduct) ValidateProduct(productId string) (dto.BaseResponse, error) {
	response, err := outbond.GetProductValidation(productId)
	if err != nil {
		fmt.Printf("Error consuming API: %v\n", err)
		return dto.BaseResponse{}, err
	}

	fmt.Printf("Product validation response: %+v\n", response)

	return *response, nil
}

func NewUsecaseProduct() UsecaseProduct {
	return &usecaseProduct{}
}
