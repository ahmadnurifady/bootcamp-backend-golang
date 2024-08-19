package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"orchestrator-order/internal/domain"
	"orchestrator-order/internal/domain/dto"
	"orchestrator-order/internal/usecase"
)

type TransactionHandler interface {
	Route()
	CreateTransactionHandler(ctx *gin.Context)
	EditRetry(ctx *gin.Context)
}

type transactionHandler struct {
	uc usecase.TransactionUsecase
	rg *gin.RouterGroup
}

func (h *transactionHandler) Route() {
	tg := h.rg.Group("/transactions")
	tg.POST("/editRetry", h.EditRetry)
}

func (h *transactionHandler) EditRetry(ctx *gin.Context) {
	var request domain.Transaction

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to parse request body",
			Data:       nil,
		})
		return
	}

	result, err := h.uc.EditRetryUsecase(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "failed to update transaction",
			Data:       nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "success edit retry",
		Data:       result,
	})
}

func (h *transactionHandler) CreateTransactionHandler(ctx *gin.Context) {
	var request domain.Transaction

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	result, err := h.uc.CreateTransaction(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	ctx.JSON(http.StatusCreated, dto.BaseResponse{
		StatusCode: http.StatusCreated,
		Message:    "Success",
		Data:       result,
	})

}

func NewTransactionHandler(uc usecase.TransactionUsecase, rg *gin.RouterGroup) TransactionHandler {
	return &transactionHandler{uc: uc, rg: rg}
}
