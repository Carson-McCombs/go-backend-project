package transaction

import (
	"time"
)

type Transaction struct {
	Payer     string    `json:"payer"`
	Points    int64     `json:"points"`
	Timestamp time.Time `json:"timestamp"`
	Buffer    int64     // Current amount of the total points a transaction has spent / withdrawn
}

// Returns a value signifying where transaction A should be place in relation to transaction B ( -1 -> before B; 0 -> equal with B; 1 -> after B)
func Compare(transactionA Transaction, transactionB Transaction) int8 {
	if transactionB.Timestamp.Before(transactionA.Timestamp) {
		return -1
	} else if transactionB.Timestamp.After(transactionA.Timestamp) {
		return 1
	} else {
		return 1
	}
}
