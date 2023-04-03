package service

import (
	"fmt"
	"mime/multipart"
	"regexp"

	domainMail "github.com/jeffleon1/transaction-ms/pkg/mail/domain"
	"github.com/jeffleon1/transaction-ms/pkg/transactions/domain"
	domainTransactions "github.com/jeffleon1/transaction-ms/pkg/transactions/domain"
)

type transactionService struct {
	transactionRepo domainTransactions.TransactionRepository
	accountRepo     domainTransactions.AccountRepository
	mailRepo        domainMail.GrpcMailRepository
}

type TransactionService interface {
	ProcessAccountData(file *multipart.File, header *multipart.FileHeader) error
}

func NewTransactionService(
	transactionRepo domainTransactions.TransactionRepository,
	accountRepo domain.AccountRepository,
	mailRepo domainMail.GrpcMailRepository,
) TransactionService {
	return &transactionService{
		transactionRepo,
		accountRepo,
		mailRepo,
	}
}

func (t *transactionService) ProcessAccountData(file *multipart.File, header *multipart.FileHeader) error {
	matched, err := regexp.MatchString(`([0-9A-Za-z]+).csv$`, header.Filename)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("not type of file allowed please try with csv file")
	}
	transactions, err := t.transactionRepo.CastMultipartFileToStruct(file)
	if err != nil {
		return err
	}

	resumeAccount := t.transactionRepo.ProccessTransactions(transactions)

	accountResume := domainTransactions.AccountResume{
		TransactionsResume: resumeAccount.Months,
		Total:              fmt.Sprintf("$%.2f", resumeAccount.Total),
		AverageCredit:      fmt.Sprintf("$%.2f", resumeAccount.Credit),
		AverageDebit:       fmt.Sprintf("$%.2f", resumeAccount.Debit),
	}

	if err := t.accountRepo.InsertAccountResume(accountResume); err != nil {
		return err
	}

	if err := t.mailRepo.SendMail(accountResume); err != nil {
		return err
	}

	return nil
}
