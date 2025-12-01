package constants

const (
	// Transaction Types
	TransactionTypeDebit  = "DEBIT"
	TransactionTypeCredit = "CREDIT"

	// Transaction Statuses
	TransactionStatusSuccess = "SUCCESS"
	TransactionStatusFailed  = "FAILED"
	TransactionStatusPending = "PENDING"
)

func IsValidTransactionType(t string) bool {
	return t == TransactionTypeDebit || t == TransactionTypeCredit
}

func IsValidTransactionStatus(t string) bool {
	return t == TransactionStatusSuccess || t == TransactionStatusFailed || t == TransactionStatusPending
}
