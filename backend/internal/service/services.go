package service

import (
	"context"

	"github.com/andibalo/flip-test/internal/entity"
)

type TransactionService interface {
	UploadCSVFile(ctx context.Context, fileContent []byte) (int, error)
	GetTotalBalance(ctx context.Context) (int64, error)
	GetUnsuccessfulTransactions(ctx context.Context, req entity.GetIssuesFilter) (*entity.IssuesResponse, int64, error)
}
