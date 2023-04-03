package infra

import (
	"context"
	"time"

	"github.com/jeffleon1/transaction-ms/internal/config"
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
		From:              config.Config.From,
		FromName:          config.Config.FromName,
		To:                mail.Email,
		Subject:           config.Config.Subject,
	})

	if err != nil {
		return err
	}

	logrus.Info(response)
	return nil
}
