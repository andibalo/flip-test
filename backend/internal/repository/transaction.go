package repository

import (
	"sort"
	"sync"

	"github.com/andibalo/flip-test/internal/entity"
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

func (r *transactionRepository) GetUnsuccessfulTransactions(filter entity.GetIssuesFilter) ([]*model.Transaction, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.transactions) == 0 {
		return []*model.Transaction{}, 0, nil
	}

	unsuccessfulTransactions := make([]*model.Transaction, 0)

	for _, tx := range r.transactions {
		if tx.Status == "FAILED" || tx.Status == "PENDING" {
			unsuccessfulTransactions = append(unsuccessfulTransactions, tx)
		}
	}

	if len(filter.Sorts.Data()) > 0 {
		sortData := filter.Sorts.Data()[0]

		sort.Slice(unsuccessfulTransactions, func(i, j int) bool {
			var less bool
			switch sortData.Name {
			case "timestamp":
				less = unsuccessfulTransactions[i].TransactionDate < unsuccessfulTransactions[j].TransactionDate
			case "name":
				less = unsuccessfulTransactions[i].Name < unsuccessfulTransactions[j].Name
			case "type":
				less = unsuccessfulTransactions[i].Type < unsuccessfulTransactions[j].Type
			case "amount":
				less = unsuccessfulTransactions[i].Amount < unsuccessfulTransactions[j].Amount
			case "status":
				less = unsuccessfulTransactions[i].Status < unsuccessfulTransactions[j].Status
			case "description":
				less = unsuccessfulTransactions[i].Description < unsuccessfulTransactions[j].Description
			default:
				less = false
			}

			if sortData.Direction == "desc" {
				return !less
			}
			return less
		})
	}

	totalCount := int64(len(unsuccessfulTransactions))

	offset := (filter.GetPageWithDefault() - 1) * filter.GetPageSizeWithDefault()
	limit := filter.GetPageSizeWithDefault()

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
