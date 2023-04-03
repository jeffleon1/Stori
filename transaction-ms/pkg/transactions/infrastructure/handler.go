package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/jeffleon1/transaction-ms/pkg/transactions/application"
)

type Response struct {
	Msg    string      `json:"msg"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Err    interface{} `json:"error"`
}

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) TransactionHandler {
	return TransactionHandler{
		transactionService,
	}
}

// Account balance processor
// @Tags Account
// @Summary account balance processor
// @Description upload one csv file with the resume of account user
// @Accept  json
// @Produce  json
// @Param file	formData file true "this is a csv test file"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /account/balance [post]
func (h *TransactionHandler) AccountBalanceProcessor(ctx *gin.Context) {
	csvPartFile, csvHeader, err := ctx.Request.FormFile("file") //env variable
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Msg:    "something goes wrong please check your csv file",
			Status: "ERROR",
			Data:   nil,
			Err:    err.Error(),
		})
		return
	}

	if err := h.transactionService.ProcessAccountData(&csvPartFile, csvHeader); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Msg:    "something goes wrong please check your csv file",
			Status: "ERROR",
			Data:   nil,
			Err:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Msg:    "your resume account was created succedded",
		Status: "OK",
		Data:   nil,
		Err:    nil,
	})
}
