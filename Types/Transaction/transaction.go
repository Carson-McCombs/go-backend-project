package transaction

import (
	"time"
)

type Transaction struct {
	Payer     string    `json:"payer"`
	Points    int64     `json:"points"`
	Timestamp time.Time `json:"timestamp"`
	Buffer    int64
}

func Compare(transactionA Transaction, transactionB Transaction) int8 {
	if transactionB.Timestamp.Before(transactionA.Timestamp) {
		return -1
	} else if transactionB.Timestamp.After(transactionA.Timestamp) {
		return 1
	} else {
		return 1
	}
}
