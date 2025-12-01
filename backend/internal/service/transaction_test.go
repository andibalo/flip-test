package service

import (
	"context"
	"testing"

	"github.com/andibalo/flip-test/internal/model"
	"github.com/andibalo/flip-test/internal/repository/mocks"
	"github.com/andibalo/flip-test/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadCSVFile(t *testing.T) {
	tests := []struct {
		name          string
		csvContent    []byte
		setupMock     func(*mocks.MockTransactionRepository)
		expectedCount int
		expectError   bool
		errorContains string
	}{
		{
			name:          "Error - Empty file",
			csvContent:    []byte(``),
			setupMock:     func(mockRepo *mocks.MockTransactionRepository) {},
			expectedCount: 0,
			expectError:   true,
			errorContains: "CSV file is empty",
		},
		{
			name:          "Error - Invalid column count (less than 6)",
			csvContent:    []byte(`1624608050, MERCHANT A, CREDIT, 500000`),
			setupMock:     func(mockRepo *mocks.MockTransactionRepository) {},
			expectedCount: 0,
			expectError:   true,
			errorContains: "invalid CSV format",
		},
		{
			name:       "Error - Invalid timestamp format",
			csvContent: []byte(`invalid_timestamp, MERCHANT A, CREDIT, 500000, SUCCESS, salary`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Maybe()
			},
			expectedCount: 0,
			expectError:   true,
			errorContains: "invalid timestamp",
		},
		{
			name:       "Error - Invalid amount format",
			csvContent: []byte(`1624608050, MERCHANT A, CREDIT, invalid_amount, SUCCESS, salary`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Maybe()
			},
			expectedCount: 0,
			expectError:   true,
			errorContains: "invalid amount",
		},
		{
			name:       "Error - Zero amount",
			csvContent: []byte(`1624608050, MERCHANT A, CREDIT, 0, SUCCESS, salary`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Maybe()
			},
			expectedCount: 0,
			expectError:   true,
			errorContains: "invalid amount",
		},
		{
			name:       "Error - Negative amount",
			csvContent: []byte(`1624608050, MERCHANT A, CREDIT, -500000, SUCCESS, salary`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Maybe()
			},
			expectedCount: 0,
			expectError:   true,
			errorContains: "invalid amount",
		},
		{
			name:       "Error - Invalid transaction type",
			csvContent: []byte(`1624608050, MERCHANT A, INVALID, 500000, SUCCESS, salary`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Maybe()
			},
			expectedCount: 0,
			expectError:   true,
			errorContains: "invalid transaction type",
		},
		{
			name:       "Error - Invalid status",
			csvContent: []byte(`1624608050, MERCHANT A, CREDIT, 500000, INVALID_STATUS, salary`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Maybe()
			},
			expectedCount: 0,
			expectError:   true,
			errorContains: "invalid status",
		},
		{
			name:       "Success - Single row",
			csvContent: []byte(`1624608050, MERCHANT A, CREDIT, 500000, SUCCESS, salary`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Once()
				mockRepo.On("Clear").Return().Maybe()
				mockRepo.On("SaveBulk", mock.AnythingOfType("[]*model.Transaction")).Return(nil).Once()
			},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "Success - All transaction types",
			csvContent: []byte(`1624608050, MERCHANT A, CREDIT, 500000, SUCCESS, salary
1624615065, MERCHANT B, DEBIT, 150000, PENDING, payment`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Once()
				mockRepo.On("Clear").Return().Maybe()
				mockRepo.On("SaveBulk", mock.AnythingOfType("[]*model.Transaction")).Return(nil).Once()
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "Success - All statuses",
			csvContent: []byte(`1624608050, MERCHANT A, CREDIT, 500000, SUCCESS, salary
1624615065, MERCHANT B, DEBIT, 150000, FAILED, payment
1624622080, MERCHANT C, CREDIT, 200000, PENDING, bonus`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Once()
				mockRepo.On("Clear").Return().Maybe()
				mockRepo.On("SaveBulk", mock.AnythingOfType("[]*model.Transaction")).Return(nil).Once()
			},
			expectedCount: 3,
			expectError:   false,
		},
		{
			name: "Success - Valid CSV with multiple rows",
			csvContent: []byte(`1624608050, MERCHANT A, CREDIT, 500000, SUCCESS, salary
1624615065, MERCHANT B, DEBIT, 150000, FAILED, payment`),
			setupMock: func(mockRepo *mocks.MockTransactionRepository) {
				mockRepo.On("GetAll").Return([]*model.Transaction{}, nil).Once()
				mockRepo.On("Clear").Return().Maybe()
				mockRepo.On("SaveBulk", mock.AnythingOfType("[]*model.Transaction")).Return(nil).Once()
			},
			expectedCount: 2,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := mocks.NewMockTransactionRepository(t)

			logger := logger.GetLoggerWithDefaultOptions()

			service := NewTransactionService(mockRepo, logger)

			tt.setupMock(mockRepo)

			count, err := service.UploadCSVFile(context.Background(), tt.csvContent)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedCount, count)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
