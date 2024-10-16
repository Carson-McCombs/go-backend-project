package account

import (
	transaction "go-fetch-backend/Types/Transaction"
	"testing"
	"time"
)

func Test_Combined(t *testing.T) {
	account := NewAccount()
	transactions := []transaction.Transaction{
		{Payer: "DANNON", Points: 300, Buffer: 300, Timestamp: time.Date(2022, 10, 31, 10, 0, 0, 0, time.UTC)},
		{Payer: "UNILEVER", Points: 200, Buffer: 200, Timestamp: time.Date(2022, 10, 31, 11, 0, 0, 0, time.UTC)},
		{Payer: "DANNON", Points: -200, Buffer: -200, Timestamp: time.Date(2022, 10, 31, 15, 0, 0, 0, time.UTC)},
		{Payer: "MILLER COORS", Points: 10000, Buffer: 10000, Timestamp: time.Date(2022, 11, 1, 14, 0, 0, 0, time.UTC)},
		{Payer: "DANNON", Points: 1000, Buffer: 1000, Timestamp: time.Date(2022, 11, 2, 14, 0, 0, 0, time.UTC)},
		{Payer: "", Points: -5000, Buffer: -5000, Timestamp: time.Now()},
	}
	for _, item := range transactions {
		if item.Points < 0 {
			account.WithdrawlTransaction(item)
		} else if item.Points > 0 {
			account.DepositTransaction(item)
		}

	}
	expectedTotalsMap := map[string]uint64{"DANNON": 1000, "UNILEVER": 0, "MILLER COORS": 5300}
	totalsMap := account.GetBalanceTotalsMap()
	for payer := range totalsMap {
		if totalsMap[payer] != expectedTotalsMap[payer] {
			t.Fatalf("Error: \n Expected total balance map: %+v \n Actual total balance map: %+v", expectedTotalsMap, totalsMap)
		}
	}
}
