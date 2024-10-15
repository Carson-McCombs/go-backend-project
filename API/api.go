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
	*mux.Router
	account *account.Account
}

func NewServer() *Server {
	server := &Server{
		Router:  mux.NewRouter(),
		account: account.NewAccount(),
	}
	server.routes()
	return server
}

// All possible ways of interacting with the server
// API Routes
func (s *Server) routes() {
	s.HandleFunc("/add", s.addPoints).Methods("POST")
	s.HandleFunc("/spend", s.spendPoints).Methods("POST")
	s.HandleFunc("/balance", s.getBalance).Methods("GET")
}

// type idResponse struct {
// 	Id string `json:"id"`
// }

func (server *Server) addPoints(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	newTransaction := transaction.Transaction{}
	err := json.NewDecoder(request.Body).Decode(&newTransaction)
	if err != nil {
		request.Response.StatusCode = 400
		return
	}
	newTransaction.Buffer = newTransaction.Points
	success := server.account.DepositTransaction(newTransaction)
	if success {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
}

type spendRequest struct {
	Points int64 `json:"points"`
}
type withdrawlOutput struct {
	Payer  string `json:"payer"`
	Points int64  `json:"points"`
}

func transactionToWithdrawl(item transaction.Transaction) withdrawlOutput {
	return withdrawlOutput{
		Payer:  item.Payer,
		Points: item.Buffer,
	}
}

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

func (server *Server) getBalance(writer http.ResponseWriter, request *http.Request) {
	balanceMap := server.account.GetBalanceTotalsMap()
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(balanceMap)
}
