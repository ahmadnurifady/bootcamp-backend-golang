package usecase

import (
	"service-user/internal/domain"
	"service-user/internal/repository"
)

type UsecaseUser interface {
	ValidateUser(userId string) (domain.User, error)
}

type usecaseUser struct {
	repo repository.RepositoryUser
}

func (u usecaseUser) ValidateUser(userId string) (domain.User, error) {
	result, err := u.repo.FindUser(userId)
	if err != nil {
		return domain.User{}, err
	}

	return result, nil
}

func NewUsecaseUser(repo repository.RepositoryUser) UsecaseUser {
	return &usecaseUser{repo: repo}
}
