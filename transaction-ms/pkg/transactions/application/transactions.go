package service

import (
	"fmt"
	"mime/multipart"
	"regexp"

	domainMail "github.com/jeffleon1/transaction-ms/pkg/mail/domain"
	domainTransactions "github.com/jeffleon1/transaction-ms/pkg/transactions/domain"
)

type transactionService struct {
	transactionRepo domainTransactions.TransactionRepository
	accountRepo     domainTransactions.AccountRepository
	mailRepo        domainMail.GrpcMailRepository
}

type TransactionService interface {
	ProcessAccountData(file *multipart.File, header *multipart.FileHeader, mail string) error
}

func NewTransactionService(
	transactionRepo domainTransactions.TransactionRepository,
	accountRepo domainTransactions.AccountRepository,
	mailRepo domainMail.GrpcMailRepository,
) TransactionService {
	return &transactionService{
		transactionRepo,
		accountRepo,
		mailRepo,
	}
}

// This function
// 1.) macth if mail string has a correct way of email
// 2.) check if file uploaded is a csv
// 3.) do a cast of file uploaded in an array of structs
// 4.) process this array struct and return and object with all transactions in correct way
// 5.) insert the account resume in mongodb
// 6.) send the email with the resume
func (t *transactionService) ProcessAccountData(file *multipart.File, header *multipart.FileHeader, mail string) error {

	mailValidation, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, mail)
	if err != nil {
		return err
	}

	if !mailValidation {
		return fmt.Errorf("email %s is not allowed", mail)
	}

	fileTypeMatched, err := regexp.MatchString(`([0-9A-Za-z]+).csv$`, header.Filename)
	if err != nil {
		return err
	}

	if !fileTypeMatched {
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
		Email:              mail,
	}

	if err := t.accountRepo.InsertAccountResume(accountResume); err != nil {
		return err
	}

	if err := t.mailRepo.SendMail(accountResume); err != nil {
		return err
	}

	return nil
}
