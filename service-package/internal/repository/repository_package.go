package repository

import (
	"errors"
	"fmt"
	"service-package/internal/domain"
)

type RepositoryPackage interface {
	FindPackage(idPackage string) (domain.Package, error)
	CreatePackage(packageRequest *domain.Package) (*domain.Package, error)
}

type repositoryPackage struct {
	db map[string]domain.Package
}

func (repo *repositoryPackage) CreatePackage(packageRequest *domain.Package) (*domain.Package, error) {
	if _, exists := repo.db[packageRequest.Id]; exists {
		return nil, errors.New("package sudah terdaftar")
	}

	for _, existingPackage := range repo.db {
		if packageRequest.PackageName == existingPackage.PackageName {
			return nil, errors.New("nama package ini sudah dipakai, silahkan gunakan nama lain")
		}
	}

	repo.db[userRequest.Id] = *userRequest
	return userRequest, nil
}

func (repo *repositoryUser) FindUser(idUser string) (domain.Package, error) {
	user, exists := repo.db[idUser]
	if !exists {
		return domain.Package{}, fmt.Errorf("user dengan id: %s tidak ditemukan", idUser)
	}
	return user, nil
}

func NewRepositoryUser() RepositoryUser {
	return &repositoryUser{db: make(map[string]domain.Package)}
}
