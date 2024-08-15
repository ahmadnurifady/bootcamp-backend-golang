package handler

//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"orchestrator-order/internal/domain/dto"
//	"orchestrator-order/internal/usecase"
//)
//
//type TransactionHandler interface {
//	CreateTransactionHandler(ctx *gin.Context)
//}
//
//type transactionHandler struct {
//	uc usecase.TransactionUsecase
//}
//
//func (h *transactionHandler) CreateTransactionHandler(ctx *gin.Context) {
//	var request dto.TransactionRequest
//
//	err := ctx.ShouldBindJSON(&request)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, dto.BaseResponse{
//			StatusCode: http.StatusBadRequest,
//			Message:    err.Error(),
//			Data:       nil,
//		})
//	}
//
//	result, err := h.uc.CreateTransaction(request, ctx)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, dto.BaseResponse{
//			StatusCode: http.StatusBadRequest,
//			Message:    err.Error(),
//			Data:       nil,
//		})
//	}
//
//	ctx.JSON(http.StatusCreated, dto.BaseResponse{
//		StatusCode: http.StatusCreated,
//		Message:    "Success",
//		Data:       result,
//	})
//
//}
//
//func NewTransactionHandler(uc usecase.TransactionUsecase) TransactionHandler {
//	return &transactionHandler{uc: uc}
//}
