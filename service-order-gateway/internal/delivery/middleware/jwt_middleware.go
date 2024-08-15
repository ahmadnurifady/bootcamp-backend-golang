package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-order-gateway/internal/utils/jwt"
	"strings"
)

type AuthMiddleware interface {
	RequireToken() gin.HandlerFunc
}

type authMiddleware struct {
	jwtService jwt.JwtToken
}

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func (a *authMiddleware) RequireToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var aH authHeader
		if err := ctx.ShouldBindHeader(&aH); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		tokenString := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", -1)
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		_, err := a.jwtService.VerifyToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		//ctx.Set(config.UserSesion, claims["userId"])
		//ctx.Set(config.RoleSesion, claims["role"])
		// fmt.Println("claims :", claims)

		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService jwt.JwtToken) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
