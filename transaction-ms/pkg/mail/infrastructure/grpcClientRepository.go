package infra

import (
	"context"
	"time"

	domainMail "github.com/jeffleon1/transaction-ms/pkg/mail/domain"
	pb "github.com/jeffleon1/transaction-ms/pkg/mail/domain/proto"
	domainTransactions "github.com/jeffleon1/transaction-ms/pkg/transactions/domain"
	"github.com/sirupsen/logrus"
)

type MailClient struct {
	client pb.MailServiceClient
}

func NewGrpcMailClient(client pb.MailServiceClient) domainMail.GrpcMailRepository {
	return &MailClient{
		client,
	}
}

func (mc *MailClient) SendMail(mail domainTransactions.AccountResume) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	response, err := mc.client.SendMail(ctx, &pb.MailRequest{
		TransactionResume: mail.TransactionsResume,
		Total:             mail.Total,
		AverageCredit:     mail.AverageCredit,
		AverageDebit:      mail.AverageDebit,
		From:              "jeffersonleon1527@gmail.com",
		FromName:          "pepito",
		To:                "example@gmail.com",
		Subject:           "Account Resume",
	})

	if err != nil {
		return err
	}

	logrus.Info(response)
	return nil
}
