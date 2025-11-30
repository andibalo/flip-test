package entity

import "github.com/andibalo/flip-test/internal/model"

type UploadCSVResponse struct {
	TransactionsUploaded int `json:"transactions_uploaded"`
}

type BalanceResponse struct {
	TotalBalance int64 `json:"total_balance"`
}

type IssuesResponse struct {
	Transactions []*model.Transaction `json:"transactions"`
}
