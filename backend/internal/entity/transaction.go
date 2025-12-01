package entity

import (
	"github.com/andibalo/flip-test/internal/model"
	"github.com/andibalo/flip-test/pkg/pagination"
	"github.com/andibalo/flip-test/pkg/sort"
)

type UploadCSVResponse struct {
	TransactionsUploaded int `json:"transactions_uploaded"`
}

type BalanceResponse struct {
	TotalBalance int64 `json:"total_balance"`
}

type GetIssuesQueryParams struct {
	Sorts string `json:"sorts" form:"sorts"`
	pagination.PaginationRequest
}

type GetIssuesFilter struct {
	Sorts sort.Sorts
	pagination.PaginationRequest
}

type IssuesResponse struct {
	Transactions []*model.Transaction `json:"transactions"`
}
