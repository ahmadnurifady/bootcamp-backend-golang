package usecase

import (
	"fmt"
	"service-order-gateway/internal/utils/jwt"
)

type AuthorizationUsecase interface {
	Login() (string, error)
}

type authorizationUsecase struct {
	jwtToken jwt.JwtToken
}

// Login implements AuthenticationUsecase.
func (uc *authorizationUsecase) Login() (string, error) {

	token, err := uc.jwtToken.GenerateToken()
	if err != nil {
		fmt.Println("error at create token", err)
		return "", err
	}

	return token, nil
}

func NewAuthorizationUsecase(jwtToken jwt.JwtToken) AuthorizationUsecase {
	return &authorizationUsecase{jwtToken: jwtToken}
}
