package account

import (
	balance "go-fetch-backend/Types/Balance"
	binarytree "go-fetch-backend/Types/BinaryTree"
	transaction "go-fetch-backend/Types/Transaction"
)

type Account struct {
	balances        map[string]*balance.Balance
	orderedBalances []string
	balance         uint64
}

func NewAccount() *Account {
	return &Account{
		balances:        make(map[string]*balance.Balance),
		orderedBalances: []string{},
		balance:         0,
	}
}

// func (account *Account) AddTransaction(newTransaction transaction.Transaction) (bool, []transaction.Transaction) {
// 	isWithdrawl := newTransaction.Points < 0
// 	if isWithdrawl && account.balance < uint64(-newTransaction.Points) {
// 		return false, []transaction.Transaction{}
// 	}
// 	isPayerSpecified := newTransaction.Payer != ""

// 	bal, hasKey := account.balances[newTransaction.Payer]
// 	if !hasKey && isPayerSpecified {
// 		account.balances[newTransaction.Payer] = balance.NewBalance()
// 		bal = account.balances[newTransaction.Payer]
// 		account.orderedBalances = append(account.orderedBalances, newTransaction.Payer)
// 	}
// 	//bal := account.balances[newTransaction.Payer]
// 	if isPayerSpecified && isWithdrawl {
// 		if bal.Total < uint64(-newTransaction.Points) {
// 			return false, []transaction.Transaction{}
// 		}
// 		account.balance -= uint64(-newTransaction.Points)
// 		updateTimestamp := bal.WithdrawlByPayer(newTransaction)

// 		//_, updateTimestamp, _ := bal.AddTransaction(newTransaction)
// 		if updateTimestamp {
// 			account.updateBalanceOrder()
// 		}
// 		return true, []transaction.Transaction{newTransaction}
// 	} else if isPayerSpecified { //payer specific - deposit
// 		account.balance += uint64(newTransaction.Points)
// 		bal.AddTransaction(newTransaction)
// 		return true, []transaction.Transaction{}
// 	} else if isWithdrawl { //generic withdrawl
// 		account.balance -= uint64(newTransaction.Points)
// 		for {
// 			bal := account.balances[account.orderedBalances[0]]
// 			if bal.Total == 0 {
// 				return false, []transaction.Transaction{}
// 			}
// 			completed, updateTimestamp, remainingTransaction := bal.AddTransaction(newTransaction)
// 			if updateTimestamp {
// 				account.updateBalanceOrder()
// 			}
// 			if completed {
// 				break
// 			}
// 			newTransaction = remainingTransaction
// 		}

//		}
//		return true, []transaction.Transaction{}
//	}
func (account *Account) WithdrawlTransaction(newTransaction transaction.Transaction) (bool, []transaction.Transaction) {
	isWithdrawl := newTransaction.Points < 0
	if !isWithdrawl || account.balance < uint64(-newTransaction.Points) {
		return false, []transaction.Transaction{}
	}
	isPayerSpecified := newTransaction.Payer != ""

	bal, hasKey := account.balances[newTransaction.Payer]
	if !hasKey && isPayerSpecified {
		account.balances[newTransaction.Payer] = balance.NewBalance()
		bal = account.balances[newTransaction.Payer]
		account.orderedBalances = append(account.orderedBalances, newTransaction.Payer)
	}

	if isPayerSpecified {
		if bal.Total < uint64(-newTransaction.Points) {
			return false, []transaction.Transaction{}
		}
		account.balance -= uint64(-newTransaction.Points)
		updateTimestamp := bal.WithdrawlByPayer(newTransaction)
		if updateTimestamp {
			account.updateBalanceOrder()
		}
		return true, []transaction.Transaction{newTransaction}
	}
	account.balance -= uint64(-newTransaction.Points)

	withdrawlMap := map[string]transaction.Transaction{}
	for {
		bal := account.balances[account.orderedBalances[0]]
		if bal.Total == 0 {
			return false, []transaction.Transaction{}
		}
		completed, updateTimestamp, spentTransaction, remainingTransaction := bal.WithdrawlTransaction(newTransaction)
		spentTransaction.Payer = account.orderedBalances[0]
		currentWithdrawl, hasKey := withdrawlMap[spentTransaction.Payer]
		if !hasKey {
			withdrawlMap[spentTransaction.Payer] = spentTransaction
		} else {
			currentWithdrawl.Buffer += spentTransaction.Buffer
			withdrawlMap[spentTransaction.Payer] = currentWithdrawl
		}

		if updateTimestamp {
			account.updateBalanceOrder()
		}
		if completed {
			break
		}
		newTransaction = remainingTransaction

	}
	withdrawlList := []transaction.Transaction{}
	for _, withdrawl := range withdrawlMap {
		withdrawlList = append(withdrawlList, withdrawl)
	}
	return true, withdrawlList
}
func (account *Account) DepositTransaction(newTransaction transaction.Transaction) bool {
	if len(newTransaction.Payer) == 0 || newTransaction.Timestamp.IsZero() {
		return false
	}
	bal, hasKey := account.balances[newTransaction.Payer]
	if !hasKey {
		account.balances[newTransaction.Payer] = balance.NewBalance()
		bal = account.balances[newTransaction.Payer]
		account.orderedBalances = append(account.orderedBalances, newTransaction.Payer)
	}
	account.balance += uint64(newTransaction.Points)
	bal.DepositTransaction(newTransaction)
	//bal.AddTransaction(newTransaction)
	return true
}

func (account *Account) updateBalanceOrder() {
	if len(account.orderedBalances) <= 1 {
		return
	}
	compare := func(payerA string, payerB string) int8 { //first compare by whether or not there is an unspent deposit, then compare by transaction ( timestamp )
		balanceA := account.balances[payerA]
		balanceB := account.balances[payerB]
		isEmptyA := len(balanceA.UnspentDeposits) == 0
		isEmptyB := len(balanceB.UnspentDeposits) == 0
		if isEmptyA && isEmptyB {
			return 0
		} else if isEmptyA {
			return -1
		} else if isEmptyB {
			return 1
		}
		transactionA := balanceA.UnspentDeposits[0]
		transactionB := balanceB.UnspentDeposits[0]
		return transaction.Compare(transactionA, transactionB)
	}
	payer := account.orderedBalances[0]
	newOrderedBalances := account.orderedBalances[1:]
	binarytree.InsertIntoSorted(&newOrderedBalances, payer, compare)
	account.orderedBalances = newOrderedBalances
}

func (account *Account) GetBalanceTotalsMap() map[string]uint64 {
	balanceTotalsMap := make(map[string]uint64)
	for payer := range account.balances {
		balanceTotalsMap[payer] = account.balances[payer].Total
	}
	return balanceTotalsMap
}
