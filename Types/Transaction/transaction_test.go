package transaction

import (
	"encoding/json"
	"testing"
	"time"
)

type marshalFromJson_Testcase struct {
	rawJson             string
	expectedTransaction Transaction
}

func Test_MarshalFromJson(t *testing.T) {
	testcases := []marshalFromJson_Testcase{
		{
			rawJson:             `{ "payer": "DANNON", "points": 300, "timestamp": "2022-10-31T10:00:00Z" }`,
			expectedTransaction: Transaction{Payer: "DANNON", Points: 300, Timestamp: time.Date(2022, 10, 31, 10, 0, 0, 0, time.UTC)},
		},
		{
			rawJson:             `{ "payer": "UNILEVER", "points": 200, "timestamp": "2022-10-31T11:00:00Z" }`,
			expectedTransaction: Transaction{Payer: "UNILEVER", Points: 200, Timestamp: time.Date(2022, 10, 31, 11, 0, 0, 0, time.UTC)},
		},
		{
			rawJson:             `{ "payer": "DANNON", "points": -200, "timestamp": "2022-10-31T15:00:00Z" }`,
			expectedTransaction: Transaction{Payer: "DANNON", Points: -200, Timestamp: time.Date(2022, 10, 31, 15, 0, 0, 0, time.UTC)},
		},
		{
			rawJson:             `{ "payer": "MILLER COORS", "points": 10000, "timestamp": "2022-11-01T14:00:00Z" }`,
			expectedTransaction: Transaction{Payer: "MILLER COORS", Points: 10000, Timestamp: time.Date(2022, 11, 1, 14, 0, 0, 0, time.UTC)},
		},
		{
			rawJson:             `{ "payer": "DANNON", "points": 1000, "timestamp": "2022-11-02T14:00:00Z" }`,
			expectedTransaction: Transaction{Payer: "DANNON", Points: 1000, Timestamp: time.Date(2022, 11, 2, 14, 0, 0, 0, time.UTC)},
		},
	}
	for _, testcase := range testcases {

		actualTransaction := Transaction{}
		err := json.Unmarshal([]byte(testcase.rawJson), &actualTransaction)
		if err != nil {
			t.Fatalf("Error\n Expectesdd: %+v \n Actual: %+v", testcase.expectedTransaction, actualTransaction)
		}
		t.Logf("GOT: %v\n", actualTransaction)

	}

}
