package api

import (
	"encoding/json"
	account "go-fetch-backend/Types/Account"
	transaction "go-fetch-backend/Types/Transaction"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router //uses Gorilla Mux for HTTP Routing
	account     *account.Account
}

// Creates a new server pointer
func NewServer() *Server {
	server := &Server{
		Router:  mux.NewRouter(),
		account: account.NewAccount(),
	}
	server.routes()
	return server
}

// API Routes
func (s *Server) routes() {
	s.HandleFunc("/add", s.addPoints).Methods("POST")
	s.HandleFunc("/spend", s.spendPoints).Methods("POST")
	s.HandleFunc("/balance", s.getBalance).Methods("GET")
}

// Reads in a json containing a transaction. Either depositing or withdrawing the specified amount of points from the balance under the specified payer's name.
// If successful, outputs status code ok ( 200 ), otherwise status code bad request ( 400 )
func (server *Server) addPoints(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	// attempts to decode the transaction from the request, returning status code 400 Bad Request if unable to do so ( such as with invalid syntax or content )
	newTransaction := transaction.Transaction{}
	err := json.NewDecoder(request.Body).Decode(&newTransaction)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// sets the buffer to be equal to the number of points ( used with keeping track how much of each transaction is processed / completed )
	newTransaction.Buffer = newTransaction.Points
	success := false
	if newTransaction.Buffer > 0 { // if the transaction's point value is positive, then treat it as a deposit
		success = server.account.DepositTransaction(newTransaction)
	} else if newTransaction.Buffer < 0 { // if the transaction's point value is negative, then treat it as a withdraw
		success, _ = server.account.WithdrawlTransaction(newTransaction)
	} else { // if the transaction is 0, don't do anything
		success = true
	}

	if success {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
}

// stucture for reading spend / withdraw requests
type spendRequest struct {
	Points int64 `json:"points"`
}

// structure for outputing the amount of spent / withdrawn points
type withdrawlOutput struct {
	Payer  string `json:"payer"`
	Points int64  `json:"points"`
}

// converts a transaction to the trimmed down output form for spent points
func transactionToWithdrawl(item transaction.Transaction) withdrawlOutput {
	return withdrawlOutput{
		Payer:  item.Payer,
		Points: item.Buffer,
	}
}

// Reads in a json containing the amount of points being spent. Converts that into a transaction and attempts to withdraw it from the account.
// Outputs a list of payers and the corresponding amount of points contributed by deposits under each of their names
func (server *Server) spendPoints(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	newSpendRequest := spendRequest{}
	err := json.NewDecoder(request.Body).Decode(&newSpendRequest)
	if err != nil {
		request.Response.StatusCode = 400
		return
	}
	newTransaction := transaction.Transaction{Payer: "", Points: -newSpendRequest.Points, Buffer: -newSpendRequest.Points, Timestamp: time.Now()}
	success, withdrawlList := server.account.WithdrawlTransaction(newTransaction)
	trimmedWithdrawls := []withdrawlOutput{}
	for _, withdrawl := range withdrawlList {
		trimmedWithdrawls = append(trimmedWithdrawls, transactionToWithdrawl(withdrawl))
	}
	if success {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(writer).Encode(trimmedWithdrawls)
}

// Outputs each balance's total and the corresponding payer's name that the balance is grouped by
func (server *Server) getBalance(writer http.ResponseWriter, request *http.Request) {
	balanceMap := server.account.GetBalanceTotalsMap()
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(balanceMap)
}
