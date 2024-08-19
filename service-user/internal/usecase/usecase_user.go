package usecase

import (
	"fmt"
	"service-user/internal/domain/dto"
	"service-user/internal/utils/outbond"
)

type UsecaseUser interface {
	ValidateUser(userId string) (dto.BaseResponse, error)
	ValidateNumber(number string) (dto.BaseResponse, error)
}

type usecaseUser struct {
	//repo repository.RepositoryUser
}

func (u usecaseUser) ValidateNumber(number string) (dto.BaseResponse, error) {
	response, err := outbond.GetNumberValidation(number)
	if err != nil {
		return dto.BaseResponse{}, err
	}

	return *response, nil
}

func (u usecaseUser) ValidateUser(userId string) (dto.BaseResponse, error) {

	response, err := outbond.GetUserValidation(userId)
	if err != nil {
		fmt.Printf("Error consuming API: %v\n", err)
		return dto.BaseResponse{}, err
	}

	return *response, nil
}

func NewUsecaseUser() UsecaseUser {
	return &usecaseUser{}
}
