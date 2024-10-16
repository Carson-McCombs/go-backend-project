package balance

import (
	transaction "go-fetch-backend/Types/Transaction"
	"testing"
	"time"
)

func Test_Combined(t *testing.T) {
	balance := NewBalance()
	expectedTotal := uint64(700)
	transactions := []transaction.Transaction{
		{Payer: "DANNON", Points: 300, Buffer: 300, Timestamp: time.Date(2022, 10, 31, 10, 0, 0, 0, time.UTC)},
		{Payer: "DANNON", Points: -200, Buffer: -200, Timestamp: time.Date(2022, 10, 31, 15, 0, 0, 0, time.UTC)},
		{Payer: "DANNON", Points: 1000, Buffer: 1000, Timestamp: time.Date(2022, 11, 2, 14, 0, 0, 0, time.UTC)},
		{Payer: "", Points: -400, Buffer: -400, Timestamp: time.Date(2022, 12, 5, 1, 0, 0, 0, time.UTC)},
	}
	for _, item := range transactions {
		if item.Points < 0 {
			balance.WithdrawlByPayer(item) //since we are testing with a single balance
		} else if item.Points > 0 {
			balance.DepositTransaction(item)
		}
	}
	if balance.Total != expectedTotal {
		t.Fatalf("Error, Expected Total: %d \n Actual Total: %d", expectedTotal, balance.Total)
	}
}
