package repository

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jeffleon1/transaction-ms/pkg/transactions/domain"
	dto "github.com/jeffleon1/transaction-ms/pkg/transactions/domain/DTO"
	"github.com/sirupsen/logrus"
)

type transactionRepository struct{}

func NewProccesorRepository() domain.TransactionRepository {
	return &transactionRepository{}
}

// recive a multipart file and cast to a csv transaction (an Array of structs)
func (c *transactionRepository) CastMultipartFileToStruct(file *multipart.File) ([]*domain.CSVTransaction, error) {
	transactions := []*domain.CSVTransaction{}
	if err := gocsv.UnmarshalMultipartFile(file, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (c *transactionRepository) ProccessTransactions(transactions []*domain.CSVTransaction) domain.Transaction {
	chunkTransactions := c.chunksTransactions(transactions, 2)
	channelTransaction := make(chan map[int64]map[int64]dto.AmountOfTransactions)
	channelAverage := make(chan domain.AverageObj)
	var postProccedTransaction []map[int64]map[int64]dto.AmountOfTransactions
	var arrayAverages []domain.AverageObj

	for i := 0; i < len(chunkTransactions); i++ {
		go c.operateCSVTransactions(chunkTransactions[i], channelTransaction, channelAverage)
		arrayAverages = append(arrayAverages, <-channelAverage)
		arrayAverages = append(arrayAverages, <-channelAverage)
		postProccedTransaction = append(postProccedTransaction, <-channelTransaction)
	}
	close(channelTransaction)
	close(channelAverage)
	concatTransactions := c.concatTransactionChunks(postProccedTransaction)
	results := c.returnResults(concatTransactions, arrayAverages)
	return results
}

func (c *transactionRepository) returnResults(
	concatTransactions map[int64]map[int64]dto.AmountOfTransactions,
	arrayAverages []domain.AverageObj,
) domain.Transaction {
	channelAverage := make(chan domain.Transaction)
	channelTransactionResults := make(chan domain.Transaction)
	go c.processAverage(arrayAverages, channelAverage)
	go c.processTransactionResults(concatTransactions, channelTransactionResults, channelAverage)
	fmt.Println(channelTransactionResults)
	return <-channelTransactionResults
}

func (c *transactionRepository) processTransactionResults(
	transactions map[int64]map[int64]dto.AmountOfTransactions,
	channelTransactionResults chan<- domain.Transaction,
	channelAverage <-chan domain.Transaction,
) {
	transactionsInMonth := make(map[string]int64, 0)
	var total float64 = 0
	for month, value := range transactions {
		totalTransactionsInMonth := int64(0)
		for _, amountOfTransactions := range value {
			totalTransactionsInMonth += amountOfTransactions.NoTransactions
			total += amountOfTransactions.Amount
		}
		monthString := time.Month(month).String()
		transactionsInMonth[monthString] = totalTransactionsInMonth
	}

	averageResult := <-channelAverage

	channelTransactionResults <- domain.Transaction{
		Months: transactionsInMonth,
		Total:  total,
		Credit: averageResult.Credit,
		Debit:  averageResult.Debit,
	}

}

func (c *transactionRepository) processAverage(arrayAverages []domain.AverageObj, channel chan<- domain.Transaction) {
	var averageCredit float64 = 0
	var averageNoCredit int64 = 0
	var averageDebit float64 = 0
	var averageNoDebit int64 = 0

	for _, average := range arrayAverages {
		if average.Type == "credit" {
			averageCredit += average.Total
			averageNoCredit += average.No
			continue
		}
		averageDebit += average.Total
		averageNoDebit += average.No
	}

	averages := domain.Transaction{
		Credit: averageCredit / float64(averageNoCredit),
		Debit:  averageDebit / float64(averageNoDebit),
	}

	channel <- averages

}

func (c *transactionRepository) concatTransactionChunks(postProcesedChunks []map[int64]map[int64]dto.AmountOfTransactions) map[int64]map[int64]dto.AmountOfTransactions {
	if len(postProcesedChunks) < 1 {
		return nil
	}
	transactionsObj := postProcesedChunks[0]
	postProcesedChunksV2 := postProcesedChunks[1:]
	for _, postProcesedChunkV2 := range postProcesedChunksV2 {
		for month, value := range postProcesedChunkV2 {
			if _, ok := transactionsObj[month]; !ok {
				transactionsObj[month] = map[int64]dto.AmountOfTransactions{}
			}

			for day := range value {
				var amount float64 = postProcesedChunkV2[month][day].Amount
				var noTransactions int64 = postProcesedChunkV2[month][day].NoTransactions

				if value, ok := transactionsObj[month][day]; ok {
					amount += value.Amount
					noTransactions += value.NoTransactions
				}

				transactionsObj[month][day] = dto.AmountOfTransactions{
					Amount:         amount,
					NoTransactions: noTransactions,
				}

			}
		}

	}

	return transactionsObj

}

func (c *transactionRepository) chunksTransactions(transactions []*domain.CSVTransaction, numberOfChunks int) [][]*domain.CSVTransaction {
	var result [][]*domain.CSVTransaction
	for i := 0; i < numberOfChunks; i++ {

		min := (i * len(transactions) / numberOfChunks)
		max := ((i + 1) * len(transactions)) / numberOfChunks

		result = append(result, transactions[min:max])

	}
	return result
}

func (c *transactionRepository) operateCSVTransactions(
	transactions []*domain.CSVTransaction,
	channelTransaction chan<- map[int64]map[int64]dto.AmountOfTransactions,
	channelAverage chan<- domain.AverageObj,
) {
	transactionsObj := map[int64]map[int64]dto.AmountOfTransactions{}
	creditObj := domain.AverageObj{}
	debitObj := domain.AverageObj{}
	for _, transaction := range transactions {
		amount, err := strconv.ParseFloat(transaction.Amount, 32)
		NoTransactions := int64(1)
		if err != nil {
			logrus.Errorf("Transaction with ID %s, can't be marshall amount %s", transaction.ID, transaction.Amount)
			logrus.Errorf("Error %s", err.Error())
			continue
		}
		month, day, err := c.CastDate(transaction.Date)
		if err != nil {
			logrus.Errorf("Transaction with ID %s, can't be cast date %s", transaction.ID, transaction.Date)
			logrus.Errorf("Error %s", err.Error())
			continue
		}

		if amount > 0 {
			creditObj.No++
			creditObj.Total += amount
		} else {
			debitObj.No++
			debitObj.Total += amount
		}

		if _, ok := transactionsObj[month]; !ok {
			transactionsObj[month] = map[int64]dto.AmountOfTransactions{}
		}

		if value, ok := transactionsObj[month][day]; ok {
			amount += value.Amount
			NoTransactions += value.NoTransactions
		}

		transactionsObj[month][day] = dto.AmountOfTransactions{
			Amount:         amount,
			NoTransactions: NoTransactions,
		}

	}

	creditObj.Type = "credit"
	debitObj.Type = "debit"
	channelAverage <- creditObj
	channelAverage <- debitObj
	channelTransaction <- transactionsObj
}

func (c *transactionRepository) CastDate(date string) (int64, int64, error) {
	split := strings.Split(date, "/")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("incorrect format for date %s", date)
	}
	month, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid format for month %s in date, must be have a number", split[0])
	}
	day, err := strconv.ParseInt(split[1], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid format for day %s in date, must be have a number", split[1])

	}

	return month, day, nil
}
