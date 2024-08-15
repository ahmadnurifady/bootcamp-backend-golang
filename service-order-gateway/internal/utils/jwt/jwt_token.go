package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"service-order-gateway/internal/config"
	"time"
)

type JwtToken interface {
	GenerateToken() (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type jwtToken struct {
	cfg config.TokenConfig
}

// GenerateToken implements JwtToken.
func (j *jwtToken) GenerateToken() (string, error) {
	claims := JwtTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.IssuerName,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: nil,
		},
	}

	jwtNewClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtNewClaims.SignedString(j.cfg.JwtSignatureKey)
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}

// VerifyToken implements JwtToken.
func (j *jwtToken) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.cfg.JwtSignatureKey, nil
	})

	if err != nil {
		return nil, errors.New("failed to verify token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok || claims["iss"] != j.cfg.IssuerName {
		return nil, errors.New("invalid claim token")
	}

	return claims, nil
}

func NewJwtToken(cfg config.TokenConfig) JwtToken {
	return &jwtToken{cfg: cfg}
}
