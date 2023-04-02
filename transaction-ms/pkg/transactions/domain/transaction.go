package domain

import (
	"mime/multipart"
)

type TransactionRepository interface {
	CastMultipartFileToStruct(file *multipart.File) ([]*CSVTransaction, error)
	ProccessTransactions(transactions []*CSVTransaction) Transaction
}

type Transaction struct {
	Months map[string]int64
	Total  float64
	Credit float64
	Debit  float64
}

type AverageObj struct {
	No    int64
	Total float64
	Type  string
}

type CSVTransaction struct {
	ID     string `csv:"Id"`
	Date   string `csv:"Date"`
	Amount string `csv:"Transaction"`
}
