package model

type Transaction struct {
	TransactionDate int64  `json:"timestamp"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Amount          int64  `json:"amount"`
	Status          string `json:"status"`
	Description     string `json:"description"`
}
