package repository

import (
	"context"

	"github.com/jeffleon1/transaction-ms/pkg/transactions/domain"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountRepository struct {
	db                *mongo.Database
	accountCollection *mongo.Collection
	ctx               context.Context
}

func NewAccountRepository(db *mongo.Database) domain.AccountRepository {
	ctx := context.Background()
	accountCollection := db.Collection("account_resume")

	return &accountRepository{
		db,
		accountCollection,
		ctx,
	}
}

func (t *accountRepository) InsertAccountResume(accountResume domain.AccountResume) error {
	result, err := t.accountCollection.InsertOne(t.ctx, accountResume)
	if err != nil {
		return err
	}
	newID := result.InsertedID
	logrus.Info("account resume inserted with ID ", newID)
	logrus.Info("account resume body: ", result)

	return nil
}
