package repository

import "github.com/andibalo/flip-test/internal/model"

//go:generate mockery --name=TransactionRepository
type TransactionRepository interface {
	SaveBulk(transactions []*model.Transaction) error
	GetAll() ([]*model.Transaction, error)
	GetUnsuccessfulTransactions(page, pageSize int) ([]*model.Transaction, int64, error)
	Clear()
}
