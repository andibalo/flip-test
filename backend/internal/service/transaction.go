package service

import (
	"context"
	"encoding/csv"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/andibalo/flip-test/internal/constants"
	"github.com/andibalo/flip-test/internal/entity"
	"github.com/andibalo/flip-test/internal/model"
	"github.com/andibalo/flip-test/internal/repository"
	"github.com/andibalo/flip-test/pkg/httpresp"
	"github.com/andibalo/flip-test/pkg/logger"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

type transactionService struct {
	transactionRepo repository.TransactionRepository
	logger          logger.Logger
}

func NewTransactionService(transactionRepo repository.TransactionRepository, logger logger.Logger) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		logger:          logger,
	}
}

func (s *transactionService) UploadCSVFile(ctx context.Context, fileContent []byte) (int, error) {
	reader := csv.NewReader(strings.NewReader(string(fileContent)))
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Failed to parse CSV file", zap.Error(err))
		return 0, oops.
			Code(httpresp.BadRequest.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
			With("error", err.Error()).
			Errorf("failed to parse CSV file")
	}

	if len(records) == 0 {
		s.logger.WarnWithContext(ctx, "[UploadCSVFile] CSV file is empty")
		return 0, oops.
			Code(httpresp.BadRequest.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
			Errorf("CSV file is empty")
	}

	transactions, err := s.parseCsvFile(ctx, records)
	if err != nil {
		s.logger.WarnWithContext(ctx, "[UploadCSVFile] Error parseCsvFile", zap.Error(err))
		return 0, err
	}

	err = s.transactionRepo.SaveBulk(transactions)
	if err != nil {
		s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Error SaveBulk", zap.Error(err))
		return 0, oops.
			Code(httpresp.ServerError.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).
			With("error", err.Error()).
			Errorf("failed to save transactions")
	}

	s.logger.InfoWithContext(ctx, "[UploadCSVFile] Successfully uploaded transactions",
		zap.Int("count", len(transactions)))

	return len(transactions), nil
}

func (s *transactionService) parseCsvFile(ctx context.Context, records [][]string) ([]*model.Transaction, error) {

	transactions := make([]*model.Transaction, 0, len(records))

	for i, record := range records {
		if len(record) != 6 {
			s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Invalid CSV format",
				zap.Int("line", i+1),
				zap.Int("columns", len(record)))
			return nil, oops.
				Code(httpresp.BadRequest.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
				With("line", i+1).
				With("columns", len(record)).
				Errorf("invalid CSV format: expected 6 columns, got %d at line %d", len(record), i+1)
		}

		timestamp, err := strconv.ParseInt(strings.TrimSpace(record[0]), 10, 64)
		if err != nil {
			s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Invalid timestamp",
				zap.Int("line", i+1),
				zap.Error(err))
			return nil, oops.
				Code(httpresp.BadRequest.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
				With("line", i+1).
				With("error", err.Error()).
				Errorf("invalid timestamp at line %d", i+1)
		}

		amount, err := strconv.ParseInt(strings.TrimSpace(record[3]), 10, 64)
		if err != nil {
			s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Invalid amount",
				zap.Int("line", i+1),
				zap.Error(err))
			return nil, oops.
				Code(httpresp.BadRequest.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
				With("line", i+1).
				With("error", err.Error()).
				Errorf("invalid amount at line %d", i+1)
		}

		if amount < 0 || amount == 0 {
			s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Invalid amount. Must be greater than 0", zap.Int64("amount", amount))
			return nil, oops.
				Code(httpresp.BadRequest.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
				With("line", i+1).
				With("amount", amount).
				Errorf("invalid amount at line %d: must be greater than 0", i+1)
		}

		txType := strings.TrimSpace(record[2])
		if !constants.IsValidTransactionType(txType) {
			s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Invalid transaction type",
				zap.Int("line", i+1),
				zap.String("type", txType))
			return nil, oops.
				Code(httpresp.BadRequest.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
				With("line", i+1).
				With("type", txType).
				Errorf("invalid transaction type at line %d: must be DEBIT or CREDIT", i+1)
		}

		status := strings.TrimSpace(record[4])
		if !constants.IsValidTransactionStatus(status) {
			s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Invalid status",
				zap.Int("line", i+1),
				zap.String("status", status))
			return nil, oops.
				Code(httpresp.BadRequest.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
				With("line", i+1).
				With("status", status).
				Errorf("invalid status at line %d: must be SUCCESS, FAILED, or PENDING", i+1)
		}

		existingTransactions, err := s.transactionRepo.GetAll()
		if err != nil {
			s.logger.ErrorWithContext(ctx, "[UploadCSVFile] Failed to retrieve existing transactions", zap.Error(err))
			return nil, oops.
				Code(httpresp.ServerError.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).
				With("error", err.Error()).
				Errorf("failed to retrieve transactions")
		}

		if len(existingTransactions) > 0 {
			s.transactionRepo.Clear()
		}

		transaction := &model.Transaction{
			TransactionDate: timestamp,
			Name:            strings.TrimSpace(record[1]),
			Type:            txType,
			Amount:          amount,
			Status:          status,
			Description:     strings.TrimSpace(record[5]),
		}

		transactions = append(transactions, transaction)
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].TransactionDate > transactions[j].TransactionDate
	})

	return transactions, nil
}

func (s *transactionService) GetTotalBalance(ctx context.Context) (int64, error) {
	transactions, err := s.transactionRepo.GetAll()
	if err != nil {
		s.logger.ErrorWithContext(ctx, "[GetTotalBalance] Failed to retrieve transactions", zap.Error(err))
		return 0, oops.
			Code(httpresp.ServerError.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).
			With("error", err.Error()).
			Errorf("failed to retrieve transactions")
	}

	var balance int64 = 0

	for _, tx := range transactions {
		if tx.Status != constants.TransactionStatusSuccess {
			continue
		}

		if tx.Type == constants.TransactionTypeCredit {
			balance += tx.Amount
		} else if tx.Type == constants.TransactionTypeDebit {
			balance -= tx.Amount
		}
	}

	return balance, nil
}

func (s *transactionService) GetUnsuccessfulTransactions(ctx context.Context, req entity.GetIssuesFilter) ([]*model.Transaction, int64, error) {
	transactions, totalCount, err := s.transactionRepo.GetUnsuccessfulTransactions(req)
	if err != nil {
		s.logger.ErrorWithContext(ctx, "[GetUnsuccessfulTransactions] Failed to retrieve unsuccessful transactions", zap.Error(err))
		return nil, 0, oops.
			Code(httpresp.ServerError.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).
			With("error", err.Error()).
			Errorf("failed to retrieve transactions")
	}

	return transactions, totalCount, nil
}
