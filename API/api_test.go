package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestServer_AddPoints(t *testing.T) {
	server := NewServer()
	url := "localhost:8000/"
	writer := httptest.NewRecorder()
	rawJson := `{
		"payer" : "DANNON",
		"points" : 5000,
		"timestamp" : "2020-11-02T14:00:00Z"
		}`
	request := httptest.NewRequest("POST", url+"add", strings.NewReader(rawJson))
	server.addPoints(writer, request)
	t.Fatalf("Status Code: %d", writer.Result().StatusCode)
}

func TestServer_SpendPoints(t *testing.T) {
	server := NewServer()
	url := "localhost:8000/"
	writer := httptest.NewRecorder()
	rawJson := `{
		"points" : 500
		}`
	request := httptest.NewRequest("POST", url+"spend", bytes.NewBufferString(rawJson))
	server.spendPoints(writer, request)
	outputMap := map[string]int64{}
	json.NewDecoder(writer.Result().Body).Decode(&outputMap)
	t.Logf("MAP: %+v", outputMap)
	t.Fatalf("Status Code: %d", writer.Result().StatusCode)
}

func TestServerA(t *testing.T) {
	server := NewServer()
	url := "localhost:8000/"

	writer := httptest.NewRecorder()
	rawJson := `{ "payer" : "DANNON", "points" : 5000, "timestamp" : "2020-11-02T14:00:00Z" }`
	request := httptest.NewRequest("POST", url+"add", strings.NewReader(rawJson))
	server.addPoints(writer, request)

	writer = httptest.NewRecorder()
	rawJson = `{ "points" : 500 }`
	request = httptest.NewRequest("POST", url+"spend", bytes.NewBufferString(rawJson))
	server.spendPoints(writer, request)
	outputMap := map[string]int64{}
	json.NewDecoder(writer.Result().Body).Decode(&outputMap)
	t.Logf("MAP: %+v", outputMap)
	t.Fatalf("Status Code: %d", writer.Result().StatusCode)
}

func TestServer_Combined(t *testing.T) {
	server := NewServer()

	testcases := []string{
		`{ "payer": "DANNON", "points": 300, "timestamp": "2022-10-31T10:00:00Z" }`,
		`{ "payer": "UNILEVER", "points": 200, "timestamp": "2022-10-31T11:00:00Z" }`,
		`{ "payer": "DANNON", "points": -200, "timestamp": "2022-10-31T15:00:00Z" }`,
		`{ "payer": "MILLER COORS", "points": 10000, "timestamp": "2022-11-01T14:00:00Z" }`,
		`{ "payer": "DANNON", "points": 1000, "timestamp": "2022-11-02T14:00:00Z" }`,
	}
	for _, rawJson := range testcases {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/add", bytes.NewBufferString(rawJson))
		server.addPoints(writer, request)
		success := writer.Result().StatusCode == 200
		if !success {
			t.Fatalf("Error adding points: %s", rawJson)
		}

	}
	rawJson := `{ "points": 5000 }`
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/spend", bytes.NewBufferString(rawJson))
	server.spendPoints(writer, request)
	withdrawls := []withdrawlOutput{}
	json.NewDecoder(writer.Result().Body).Decode(&withdrawls)
	t.Logf("Withdrawls: %+v", withdrawls)
	t.Logf("Status Code: %d", writer.Result().StatusCode)
	expectedWithdrawls := []withdrawlOutput{
		{Payer: "DANNON", Points: -100},
		{Payer: "UNILEVER", Points: -200},
		{Payer: "MILLER COORS", Points: -4700},
	}
	for _, expected := range expectedWithdrawls {
		if !slices.Contains(withdrawls, expected) {
			t.Fatalf("Error: \n Expected Withdrawls: %+v \n Actual Withdrawls: %+v", expectedWithdrawls, withdrawls)
		}
	}
	expectedBalanceMap := map[string]uint64{
		"DANNON":       1000,
		"UNILEVER":     0,
		"MILLER COORS": 5300,
	}
	writer = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/balance", nil)
	server.getBalance(writer, request)
	actualBalanceMap := map[string]uint64{}
	json.NewDecoder(writer.Result().Body).Decode(&actualBalanceMap)
	for payer, points := range actualBalanceMap {
		if expectedBalanceMap[payer] != points {
			t.Fatalf("Error \n Expected Balance Totals: %+v \n Actual Balance Totals: %+v", expectedBalanceMap, actualBalanceMap)
		}
	}
}
