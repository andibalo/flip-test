package service

// import (
// 	"context"
// 	"testing"

// 	"github.com/andibalo/flip-test/internal/constants"
// 	"github.com/andibalo/flip-test/internal/model"
// 	"github.com/andibalo/flip-test/internal/repository"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockTransactionRepository is a mock implementation of TransactionRepository
// type MockTransactionRepository struct {
// 	mock.Mock
// }

// func (m *MockTransactionRepository) SaveBulk(transactions []*model.Transaction) error {
// 	args := m.Called(transactions)
// 	return args.Error(0)
// }

// func (m *MockTransactionRepository) GetAll() ([]*model.Transaction, error) {
// 	args := m.Called()
// 	return args.Get(0).([]*model.Transaction), args.Error(1)
// }

// func (m *MockTransactionRepository) GetUnsuccessfulTransactions(page, pageSize int) ([]*model.Transaction, int64, error) {
// 	args := m.Called(page, pageSize)
// 	return args.Get(0).([]*model.Transaction), args.Get(1).(int64), args.Error(2)
// }

// func (m *MockTransactionRepository) Clear() {
// 	m.Called()
// }

// func TestUploadCSVFile_Success(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	csvContent := []byte(`1624608050, MERCHANT A, CREDIT, 500000, SUCCESS, salary
// 1624615065, MERCHANT B, DEBIT, 150000, FAILED, payment`)

// 	mockRepo.On("SaveBulk", mock.AnythingOfType("[]*model.Transaction")).Return(nil)
// 	mockLogger.On("InfoWithContext", mock.Anything, mock.Anything, mock.Anything).Return()

// 	count, err := service.UploadCSVFile(context.Background(), csvContent)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, count)
// 	mockRepo.AssertExpectations(t)
// }

// func TestUploadCSVFile_EmptyFile(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	csvContent := []byte(``)

// 	mockLogger.On("WarnWithContext", mock.Anything, mock.Anything).Return()

// 	count, err := service.UploadCSVFile(context.Background(), csvContent)

// 	assert.Error(t, err)
// 	assert.Equal(t, 0, count)
// 	assert.Contains(t, err.Error(), "CSV file is empty")
// }

// func TestUploadCSVFile_InvalidColumnCount(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	csvContent := []byte(`1624608050, MERCHANT A, CREDIT, 500000`)

// 	mockLogger.On("ErrorWithContext", mock.Anything, mock.Anything, mock.Anything).Return()

// 	count, err := service.UploadCSVFile(context.Background(), csvContent)

// 	assert.Error(t, err)
// 	assert.Equal(t, 0, count)
// 	assert.Contains(t, err.Error(), "invalid CSV format")
// }

// func TestUploadCSVFile_InvalidTimestamp(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	csvContent := []byte(`invalid_timestamp, MERCHANT A, CREDIT, 500000, SUCCESS, salary`)

// 	mockLogger.On("ErrorWithContext", mock.Anything, mock.Anything, mock.Anything).Return()

// 	count, err := service.UploadCSVFile(context.Background(), csvContent)

// 	assert.Error(t, err)
// 	assert.Equal(t, 0, count)
// 	assert.Contains(t, err.Error(), "invalid timestamp")
// }

// func TestUploadCSVFile_InvalidAmount(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	csvContent := []byte(`1624608050, MERCHANT A, CREDIT, invalid_amount, SUCCESS, salary`)

// 	mockLogger.On("ErrorWithContext", mock.Anything, mock.Anything, mock.Anything).Return()

// 	count, err := service.UploadCSVFile(context.Background(), csvContent)

// 	assert.Error(t, err)
// 	assert.Equal(t, 0, count)
// 	assert.Contains(t, err.Error(), "invalid amount")
// }

// func TestUploadCSVFile_InvalidTransactionType(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	csvContent := []byte(`1624608050, MERCHANT A, INVALID, 500000, SUCCESS, salary`)

// 	mockLogger.On("ErrorWithContext", mock.Anything, mock.Anything, mock.Anything).Return()

// 	count, err := service.UploadCSVFile(context.Background(), csvContent)

// 	assert.Error(t, err)
// 	assert.Equal(t, 0, count)
// 	assert.Contains(t, err.Error(), "invalid transaction type")
// }

// func TestUploadCSVFile_InvalidStatus(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	csvContent := []byte(`1624608050, MERCHANT A, CREDIT, 500000, INVALID_STATUS, salary`)

// 	mockLogger.On("ErrorWithContext", mock.Anything, mock.Anything, mock.Anything).Return()

// 	count, err := service.UploadCSVFile(context.Background(), csvContent)

// 	assert.Error(t, err)
// 	assert.Equal(t, 0, count)
// 	assert.Contains(t, err.Error(), "invalid status")
// }

// func TestGetTotalBalance_Success(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	transactions := []*model.Transaction{
// 		{TransactionDate: 1624608050, Name: "A", Type: constants.TransactionTypeCredit, Amount: 500000, Status: constants.TransactionStatusSuccess},
// 		{TransactionDate: 1624615065, Name: "B", Type: constants.TransactionTypeDebit, Amount: 150000, Status: constants.TransactionStatusSuccess},
// 		{TransactionDate: 1624622080, Name: "C", Type: constants.TransactionTypeCredit, Amount: 200000, Status: constants.TransactionStatusFailed},
// 		{TransactionDate: 1624629095, Name: "D", Type: constants.TransactionTypeDebit, Amount: 100000, Status: constants.TransactionStatusPending},
// 	}

// 	mockRepo.On("GetAll").Return(transactions, nil)

// 	balance, err := service.GetTotalBalance(context.Background())

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(350000), balance) // 500000 - 150000 = 350000 (only SUCCESS)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetTotalBalance_OnlyCredits(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	transactions := []*model.Transaction{
// 		{TransactionDate: 1624608050, Name: "A", Type: constants.TransactionTypeCredit, Amount: 500000, Status: constants.TransactionStatusSuccess},
// 		{TransactionDate: 1624615065, Name: "B", Type: constants.TransactionTypeCredit, Amount: 300000, Status: constants.TransactionStatusSuccess},
// 	}

// 	mockRepo.On("GetAll").Return(transactions, nil)

// 	balance, err := service.GetTotalBalance(context.Background())

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(800000), balance)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetTotalBalance_OnlyDebits(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	transactions := []*model.Transaction{
// 		{TransactionDate: 1624608050, Name: "A", Type: constants.TransactionTypeDebit, Amount: 200000, Status: constants.TransactionStatusSuccess},
// 		{TransactionDate: 1624615065, Name: "B", Type: constants.TransactionTypeDebit, Amount: 100000, Status: constants.TransactionStatusSuccess},
// 	}

// 	mockRepo.On("GetAll").Return(transactions, nil)

// 	balance, err := service.GetTotalBalance(context.Background())

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(-300000), balance)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetTotalBalance_EmptyTransactions(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	transactions := []*model.Transaction{}

// 	mockRepo.On("GetAll").Return(transactions, nil)

// 	balance, err := service.GetTotalBalance(context.Background())

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(0), balance)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetUnsuccessfulTransactions_Success(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	unsuccessfulTxns := []*model.Transaction{
// 		{TransactionDate: 1624608050, Name: "A", Type: constants.TransactionTypeCredit, Amount: 500000, Status: constants.TransactionStatusFailed},
// 		{TransactionDate: 1624615065, Name: "B", Type: constants.TransactionTypeDebit, Amount: 150000, Status: constants.TransactionStatusPending},
// 	}

// 	mockRepo.On("GetUnsuccessfulTransactions", 1, 10).Return(unsuccessfulTxns, int64(2), nil)

// 	transactions, totalCount, err := service.GetUnsuccessfulTransactions(context.Background(), 1, 10)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, len(transactions))
// 	assert.Equal(t, int64(2), totalCount)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetUnsuccessfulTransactions_Pagination(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	paginatedTxns := []*model.Transaction{
// 		{TransactionDate: 1624608050, Name: "A", Type: constants.TransactionTypeCredit, Amount: 500000, Status: constants.TransactionStatusFailed},
// 	}

// 	mockRepo.On("GetUnsuccessfulTransactions", 1, 1).Return(paginatedTxns, int64(5), nil)

// 	transactions, totalCount, err := service.GetUnsuccessfulTransactions(context.Background(), 1, 1)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, len(transactions))
// 	assert.Equal(t, int64(5), totalCount)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetUnsuccessfulTransactions_EmptyResult(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(mockRepo, mockLogger)

// 	emptyTxns := []*model.Transaction{}

// 	mockRepo.On("GetUnsuccessfulTransactions", 1, 10).Return(emptyTxns, int64(0), nil)

// 	transactions, totalCount, err := service.GetUnsuccessfulTransactions(context.Background(), 1, 10)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 0, len(transactions))
// 	assert.Equal(t, int64(0), totalCount)
// 	mockRepo.AssertExpectations(t)
// }

// // Integration-like test using real repository
// func TestTransactionService_Integration(t *testing.T) {
// 	repo := repository.NewTransactionRepository()
// 	mockLogger := new(MockLogger)
// 	service := NewTransactionService(repo, mockLogger)

// 	// Upload CSV
// 	csvContent := []byte(`1624608050, MERCHANT A, CREDIT, 500000, SUCCESS, salary
// 1624615065, MERCHANT B, DEBIT, 150000, SUCCESS, payment
// 1624622080, MERCHANT C, CREDIT, 200000, FAILED, refund
// 1624629095, MERCHANT D, DEBIT, 100000, PENDING, subscription`)

// 	mockLogger.On("InfoWithContext", mock.Anything, mock.Anything, mock.Anything).Return()
// 	mockLogger.On("ErrorWithContext", mock.Anything, mock.Anything, mock.Anything).Return()

// 	count, err := service.UploadCSVFile(context.Background(), csvContent)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 4, count)

// 	// Get total balance
// 	balance, err := service.GetTotalBalance(context.Background())
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(350000), balance) // 500000 - 150000 = 350000

// 	// Get unsuccessful transactions
// 	unsuccessfulTxns, totalCount, err := service.GetUnsuccessfulTransactions(context.Background(), 1, 10)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, len(unsuccessfulTxns))
// 	assert.Equal(t, int64(2), totalCount)
// }
