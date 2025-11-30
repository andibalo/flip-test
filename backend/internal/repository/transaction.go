package repository

import (
	"sync"

	"github.com/andibalo/flip-test/internal/model"
)

type transactionRepository struct {
	mu           sync.RWMutex
	transactions []*model.Transaction
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{
		transactions: make([]*model.Transaction, 0),
	}
}

func (r *transactionRepository) SaveBulk(transactions []*model.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.transactions = transactions

	return nil
}

func (r *transactionRepository) GetAll() ([]*model.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.transactions, nil
}

func (r *transactionRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.transactions = make([]*model.Transaction, 0)
}

func (r *transactionRepository) GetUnsuccessfulTransactions(page, pageSize int) ([]*model.Transaction, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	unsuccessfulTransactions := make([]*model.Transaction, 0)
	for _, tx := range r.transactions {
		if tx.Status == "FAILED" || tx.Status == "PENDING" {
			unsuccessfulTransactions = append(unsuccessfulTransactions, tx)
		}
	}

	totalCount := int64(len(unsuccessfulTransactions))

	offset := (page - 1) * pageSize
	limit := pageSize

	if offset >= len(unsuccessfulTransactions) {
		return []*model.Transaction{}, totalCount, nil
	}

	end := offset + limit
	if end > len(unsuccessfulTransactions) {
		end = len(unsuccessfulTransactions)
	}

	paginatedTransactions := unsuccessfulTransactions[offset:end]

	return paginatedTransactions, totalCount, nil
}
