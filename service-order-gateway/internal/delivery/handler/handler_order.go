package handler

import (
	"github.com/benebobaa/valo"
	"github.com/gin-gonic/gin"
	"net/http"
	"service-order-gateway/internal/delivery/middleware"
	"service-order-gateway/internal/domain/dto"
	"service-order-gateway/internal/usecase"
)

type HandlerOrder interface {
	CreateOrderHandler(ctx *gin.Context)
	Route()
}

type handlerOrder struct {
	uc            usecase.UsecaseOrder
	rg            *gin.RouterGroup
	jwtMiddleware middleware.AuthMiddleware
}

func (h *handlerOrder) CreateOrderHandler(ctx *gin.Context) {
	var requestOrder dto.OrderRequest

	if err := ctx.ShouldBind(&requestOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       err.Error(),
		})
		return
	}

	if err := valo.Validate(requestOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Validation failed",
			Data:       err.Error(),
		})
		return
	}

	result, err := h.uc.CreateOrderUsecase(requestOrder, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create order",
			Data:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.BaseResponse{
		StatusCode: http.StatusCreated,
		Message:    "Order created successfully",
		Data:       result,
	})
}

func (h *handlerOrder) Route() {
	og := h.rg.Group("/orders")
	og.POST("/create", h.CreateOrderHandler)
}

func NewHandlerOrder(uc usecase.UsecaseOrder, rg *gin.RouterGroup, jwtMiddleware middleware.AuthMiddleware) HandlerOrder {
	return &handlerOrder{
		uc:            uc,
		rg:            rg,
		jwtMiddleware: jwtMiddleware,
	}
}
