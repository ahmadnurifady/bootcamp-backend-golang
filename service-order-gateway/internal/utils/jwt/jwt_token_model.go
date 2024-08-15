package jwt

import "github.com/golang-jwt/jwt/v5"

type JwtTokenClaims struct {
	jwt.RegisteredClaims
	Token string
}
