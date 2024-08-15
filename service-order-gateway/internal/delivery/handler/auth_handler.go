package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-order-gateway/internal/domain/dto"
	"service-order-gateway/internal/usecase"
	"service-order-gateway/internal/utils/jwt"
)

type HandlerAuth interface {
	LoginHandler(ctx *gin.Context)
	RouteAuth()
}

type handlerAuth struct {
	uc  usecase.AuthorizationUsecase
	jwt jwt.JwtToken
	rg  *gin.RouterGroup
}

func (h *handlerAuth) RouteAuth() {
	hg := h.rg.Group("/auth")

	hg.POST("/login", h.LoginHandler)
}

func (h *handlerAuth) LoginHandler(ctx *gin.Context) {
	token, err := h.uc.Login()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       token,
	})
}

func NewHandlerAuth(jwt jwt.JwtToken, rg *gin.RouterGroup) HandlerAuth {
	return &handlerAuth{jwt: jwt, rg: rg}
}
