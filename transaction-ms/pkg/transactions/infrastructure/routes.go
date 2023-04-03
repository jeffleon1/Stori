package infrastructure

import "github.com/gin-gonic/gin"

type TransactionRoutes struct {
	transactionHandler TransactionHandler
}

func (ro *TransactionRoutes) PublicRoutes(public *gin.RouterGroup) {
	public.POST("/account/balance", ro.transactionHandler.AccountBalanceProcessor)
}

func NewRoutes(transactionHandler TransactionHandler) *TransactionRoutes {
	return &TransactionRoutes{
		transactionHandler,
	}
}
