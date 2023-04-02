package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/jeffleon1/transaction-ms/pkg/transactions/application"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) TransactionHandler {
	return TransactionHandler{
		transactionService,
	}
}

func (h *TransactionHandler) ProcessAccountResume(ctx *gin.Context) {
	csvPartFile, csvHeader, openErr := ctx.Request.FormFile("file")
	if openErr != nil {
		ctx.JSON(http.StatusBadRequest, openErr.Error())
	}
	h.transactionService.ProcessAccountData(&csvPartFile, csvHeader)

	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", csvPartFile))
}
