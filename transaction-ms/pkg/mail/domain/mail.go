package domain

import "github.com/jeffleon1/transaction-ms/pkg/transactions/domain"

type GrpcMailRepository interface {
	SendMail(mail domain.AccountResume) error
}
