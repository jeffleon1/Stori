package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountResume struct {
	ID                 *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	TransactionsResume map[string]int64    `json:"transactions_resume" bson:"transactions_resume"`
	Total              string              `json:"total" bson:"total"`
	AverageCredit      string              `json:"average_credit" bson:"average_credit"`
	AverageDebit       string              `json:"average_debit" bson:"average_debit"`
	Email              string              `json:"email" bson:"email"`
}

type AccountRepository interface {
	InsertAccountResume(accountResume AccountResume) error
}
