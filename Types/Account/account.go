package account

import (
	balance "go-fetch-backend/Types/Balance"
	binarytree "go-fetch-backend/Types/BinaryTree"
	transaction "go-fetch-backend/Types/Transaction"
)

type Account struct {
	balances        map[string]*balance.Balance //maps payers' names to their corresponding balance for the account
	orderedBalances []string                    //stores each payer's name in order by oldest unspent transaction
	balance         uint64                      //total balance for the account ( cannot be negative, therefore an unsigned variable is used )
}

// Create a new account pointer
func NewAccount() *Account {
	return &Account{
		balances:        make(map[string]*balance.Balance),
		orderedBalances: []string{},
		balance:         0,
	}
}

// Determines whether or not the transaction is general or targing a specific payer.
// Returns whether or not the withdrawl was completed ( i.e. there was enough of a balance to cover it ) and a list of withdrawl transactions showing what amount of points was withdrawn from which payer
func (account *Account) WithdrawlTransaction(newTransaction transaction.Transaction) (bool, []transaction.Transaction) {
	isWithdrawl := newTransaction.Points < 0
	if !isWithdrawl || account.balance < uint64(-newTransaction.Points) {
		return false, []transaction.Transaction{}
	}
	isPayerSpecified := newTransaction.Payer != ""

	bal, hasKey := account.balances[newTransaction.Payer] //create a new balance for the corresponding payer if one does not already exist
	if !hasKey && isPayerSpecified {
		account.balances[newTransaction.Payer] = balance.NewBalance()
		bal = account.balances[newTransaction.Payer]
		account.orderedBalances = append(account.orderedBalances, newTransaction.Payer)
	}

	if isPayerSpecified { //see if the withdrawl is targeting a specific payer
		if bal.Total < uint64(-newTransaction.Points) {
			return false, []transaction.Transaction{} //if that payer doesn't have enough points, return
		}
		account.balance -= uint64(-newTransaction.Points)
		updateTimestamp := bal.WithdrawlByPayer(newTransaction) //if it is, withdraw from that payer, updating the balance order if necessary
		if updateTimestamp {
			account.updateBalanceOrder()
		}
		return true, []transaction.Transaction{newTransaction}
	}
	//subtract the withdrawl from the account's overall balance if it is a generic withdrawl
	account.balance -= uint64(-newTransaction.Points)

	//loop through each payer, starting from the one with the oldest / top-priority unspent deposits / transactions
	withdrawlMap := map[string]transaction.Transaction{}
	for {
		bal := account.balances[account.orderedBalances[0]]
		if bal.Total == 0 { // if the balance is empty exit ( should never occur )
			return false, []transaction.Transaction{}
		}

		completed, updateTimestamp, spentTransaction, remainingTransaction := bal.WithdrawlTransaction(newTransaction)
		spentTransaction.Payer = account.orderedBalances[0]              //update the spent transaction with the payer name whose points were used
		currentWithdrawl, hasKey := withdrawlMap[spentTransaction.Payer] //get the transaction entry from the withdrawl map
		if !hasKey {                                                     //if one doesn't exist, set the new spent transaction as the total withdrawl for that payer
			withdrawlMap[spentTransaction.Payer] = spentTransaction
		} else { //otherwise, add the unspent buffer ( the amount spent by that payer ) to the total withdrawl transaction for that payer
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

// deposits a set amount of points into a balance with the corresponding payer ( creating a new balance if one doesn't already exist )
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
	return true
}

// updates which balance has the oldest / top-priority unspent transaction / deposit
func (account *Account) updateBalanceOrder() {
	if len(account.orderedBalances) <= 1 {
		return
	}
	//first compare by whether or not there is an unspent deposit, then compare by transaction ( timestamp )
	compare := func(payerA string, payerB string) int8 {
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

// returns a map of all of the account balances grouped by payer
func (account *Account) GetBalanceTotalsMap() map[string]uint64 {
	balanceTotalsMap := make(map[string]uint64)
	for payer := range account.balances {
		balanceTotalsMap[payer] = account.balances[payer].Total
	}
	return balanceTotalsMap
}
