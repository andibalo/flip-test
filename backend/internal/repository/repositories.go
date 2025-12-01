package repository

import (
	"github.com/andibalo/flip-test/internal/entity"
	"github.com/andibalo/flip-test/internal/model"
)

type TransactionRepository interface {
	SaveBulk(transactions []*model.Transaction) error
	GetAll() ([]*model.Transaction, error)
	GetUnsuccessfulTransactions(filter entity.GetIssuesFilter) ([]*model.Transaction, int64, error)
	GetUnsuccessfulTransactionsSummary(filter entity.GetIssuesFilter) (int64, int64, int64, error)
	Clear()
}
