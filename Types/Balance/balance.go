package balance

import (
	binarytree "go-fetch-backend/Types/BinaryTree"
	transaction "go-fetch-backend/Types/Transaction"
)

type Balance struct {
	UnspentDeposits []transaction.Transaction
	spentDeposits   []transaction.Transaction
	withdrawls      []transaction.Transaction
	Total           uint64
}

func NewBalance() *Balance {
	return &Balance{
		UnspentDeposits: []transaction.Transaction{},
		spentDeposits:   []transaction.Transaction{},
		withdrawls:      []transaction.Transaction{},
		Total:           0,
	}
}

func (balance *Balance) WithdrawlTransaction(withdrawl transaction.Transaction) (bool, bool, transaction.Transaction, transaction.Transaction) {
	completedTransaction := false
	updateBalanceOrder := false
	depositBuffer := balance.UnspentDeposits[0].Buffer
	if depositBuffer < -withdrawl.Buffer { //can't fulfill entire transaction
		withdrawl.Buffer += balance.UnspentDeposits[0].Buffer
		balance.UnspentDeposits[0].Buffer = 0
		balance.spentDeposits = append(balance.spentDeposits, binarytree.RemoveAtIndex(&balance.UnspentDeposits, 0))
		balance.Total -= uint64(depositBuffer)
		spentTransaction := withdrawl
		spentTransaction.Buffer = -depositBuffer
		balance.withdrawls = append(balance.withdrawls, spentTransaction)
		updateBalanceOrder = true
		return completedTransaction, updateBalanceOrder, spentTransaction, withdrawl
	} else if depositBuffer == -withdrawl.Buffer {
		balance.UnspentDeposits[0].Buffer = 0
		balance.spentDeposits = append(balance.spentDeposits, balance.UnspentDeposits[0])
		balance.withdrawls = append(balance.withdrawls, withdrawl)
		balance.Total -= uint64(withdrawl.Points)
		completedTransaction = true
		updateBalanceOrder = true
		return completedTransaction, updateBalanceOrder, withdrawl, transaction.Transaction{}
	}
	balance.UnspentDeposits[0].Buffer += withdrawl.Buffer
	balance.Total -= uint64(-withdrawl.Buffer)
	balance.withdrawls = append(balance.withdrawls, withdrawl)
	completedTransaction = true
	return completedTransaction, updateBalanceOrder, withdrawl, transaction.Transaction{}
}

func (balance *Balance) DepositTransaction(deposit transaction.Transaction) bool {
	updateTimestamp := len(balance.UnspentDeposits) == 0 || deposit.Timestamp.Before(balance.UnspentDeposits[0].Timestamp)
	balance.Total += uint64(deposit.Points)
	binarytree.InsertIntoSorted(&balance.UnspentDeposits, deposit, transaction.Compare)
	return updateTimestamp
}

// Will spend the remaining deposit to complete the transaction. If there is still a remaining charge, then return the remaining charge as a new transaction.
// func (balance *Balance) AddTransaction(item transaction.Transaction) (bool, bool, transaction.Transaction) {
// 	completedTransaction := false
// 	updateBalanceOrder := false
// 	if item.Points < 0 { //withdrawl

// 		depositBuffer := balance.UnspentDeposits[0].Buffer
// 		if depositBuffer < -item.Buffer { //can't fulfill entire transaction
// 			item.Buffer += balance.UnspentDeposits[0].Buffer
// 			balance.UnspentDeposits[0].Buffer = 0
// 			balance.spentDeposits = append(balance.spentDeposits, binarytree.RemoveAtIndex(&balance.UnspentDeposits, 0))
// 			balance.Total -= uint64(depositBuffer)
// 			partialItem := item
// 			partialItem.Buffer = -depositBuffer
// 			balance.withdrawls = append(balance.withdrawls, partialItem)
// 			updateBalanceOrder = true
// 			return completedTransaction, updateBalanceOrder, item
// 		} else if depositBuffer == -item.Buffer {
// 			balance.UnspentDeposits[0].Buffer = 0
// 			balance.spentDeposits = append(balance.spentDeposits, balance.UnspentDeposits[0])
// 			balance.withdrawls = append(balance.withdrawls, item)
// 			balance.Total -= uint64(item.Points)
// 			completedTransaction = true
// 			updateBalanceOrder = true
// 			return completedTransaction, updateBalanceOrder, transaction.Transaction{}
// 		}
// 		balance.UnspentDeposits[0].Buffer += item.Buffer
// 		balance.Total -= uint64(-item.Buffer)
// 		balance.withdrawls = append(balance.withdrawls, item)
// 		updateBalanceOrder = false
// 		return completedTransaction, updateBalanceOrder, transaction.Transaction{}
// 	}
// 	updateTimestamp := len(balance.UnspentDeposits) == 0 || item.Timestamp.Before(balance.UnspentDeposits[0].Timestamp)
// 	balance.Total += uint64(item.Points)
// 	binarytree.InsertIntoSorted(&balance.UnspentDeposits, item, transaction.Compare)
// 	return true, updateTimestamp, transaction.Transaction{}

// }

func (balance *Balance) WithdrawlByPayer(item transaction.Transaction) bool {
	updateBalanceOrder := balance.UnspentDeposits[0].Buffer <= item.Buffer
	balance.Total -= uint64(-item.Buffer)
	for {
		if balance.UnspentDeposits[0].Buffer >= -item.Buffer {
			balance.UnspentDeposits[0].Buffer += item.Buffer
			balance.withdrawls = append(balance.withdrawls, item)
			return updateBalanceOrder
		}
		item.Buffer -= balance.UnspentDeposits[0].Buffer
		balance.spentDeposits = append(balance.spentDeposits, binarytree.RemoveAtIndex(&balance.UnspentDeposits, 0))

	}

}
