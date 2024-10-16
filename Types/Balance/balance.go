package balance

import (
	binarytree "go-fetch-backend/Types/BinaryTree"
	transaction "go-fetch-backend/Types/Transaction"
)

type Balance struct {
	UnspentDeposits []transaction.Transaction //all transactions that deposit / add points to the account and have not been completely spent - ordered by timestamp from oldest to newest ( oldest at top )
	spentDeposits   []transaction.Transaction //all transactions that deposit / add points to the account and have been completely spent
	withdrawls      []transaction.Transaction //all transactions that withdraw / subtract points from the account
	Total           uint64                    // total balance between an account and a specifc payer
}

// Creates a new balance pointer
func NewBalance() *Balance {
	return &Balance{
		UnspentDeposits: []transaction.Transaction{},
		spentDeposits:   []transaction.Transaction{},
		withdrawls:      []transaction.Transaction{},
		Total:           0,
	}
}

// Withdraws points from the account's balance ( under a specified payer's name ) until the top / oldest unspent deposit / transaction is used
// returns: completes transaction, update balance order, spent transaction ( new transaction created that records the amount of points withdrawn under the payer's namee ), and the remaining transaction ( how much is left in the withdraw after the spent transaction )
// transactions must be a negative Point / Buffer value ( zeros are ignored )
func (balance *Balance) WithdrawlTransaction(withdrawl transaction.Transaction) (bool, bool, transaction.Transaction, transaction.Transaction) {
	if withdrawl.Buffer > 0 {
		return true, false, transaction.Transaction{}, transaction.Transaction{}
	}
	completedTransaction := false
	updateBalanceOrder := false
	depositBuffer := balance.UnspentDeposits[0].Buffer
	if depositBuffer < -withdrawl.Buffer { //if the transaction can't be fulfilled completely, fill what is possible and return the leftover as a new transaction. And refresh the update order
		withdrawl.Buffer += balance.UnspentDeposits[0].Buffer
		balance.UnspentDeposits[0].Buffer = 0
		balance.spentDeposits = append(balance.spentDeposits, binarytree.RemoveAtIndex(&balance.UnspentDeposits, 0))
		balance.Total -= uint64(depositBuffer)
		spentTransaction := withdrawl
		spentTransaction.Buffer = -depositBuffer
		balance.withdrawls = append(balance.withdrawls, spentTransaction)
		updateBalanceOrder = true
		return completedTransaction, updateBalanceOrder, spentTransaction, withdrawl
	} else if depositBuffer == -withdrawl.Buffer { //if the transaction is exactly equal to the leftover buffer on the top transaction, return an empty transaction, mark the withdrawl as complete, and refresh the update order.
		balance.UnspentDeposits[0].Buffer = 0
		balance.spentDeposits = append(balance.spentDeposits, balance.UnspentDeposits[0])
		balance.withdrawls = append(balance.withdrawls, withdrawl)
		balance.Total -= uint64(withdrawl.Points)
		completedTransaction = true
		updateBalanceOrder = true
		return completedTransaction, updateBalanceOrder, withdrawl, transaction.Transaction{}
	}
	//if the transaction can be fulfilled without completely spending the top transaction, then mark the withdrawl as complete and
	balance.UnspentDeposits[0].Buffer += withdrawl.Buffer
	balance.Total -= uint64(-withdrawl.Buffer)
	balance.withdrawls = append(balance.withdrawls, withdrawl)
	completedTransaction = true
	return completedTransaction, updateBalanceOrder, withdrawl, transaction.Transaction{}
}

// Deposits points from the account's balance ( under a specified payer's name )
// returns: whether or not the balance order needs to be updated
// transaction must have a positive Points / Buffer value ( zeros are ignored )
func (balance *Balance) DepositTransaction(deposit transaction.Transaction) bool {
	if deposit.Buffer < 0 {
		return false
	}
	balance.Total += uint64(deposit.Points)
	if len(balance.UnspentDeposits) == 0 || deposit.Timestamp.Before(balance.UnspentDeposits[0].Timestamp) { //if the new deposit has a timestamp before the top unspent transaction ( or there are no unspent transactions ), then set the new deposit on top and update the balance order
		binarytree.Insert(&balance.UnspentDeposits, 0, deposit)
		return true
	}
	//otherwise insert the deposit into the sorted unspent deposits
	binarytree.InsertIntoSorted(&balance.UnspentDeposits, deposit, transaction.Compare)
	return false
}

// Withdraws points from the account's balance ( under a specified payer's name ) until the withdraw is completed
// returns: whether or not the balance order needs to be updated
// transactions must be a negative Point / Buffer value ( zeros are ignored )
func (balance *Balance) WithdrawlByPayer(item transaction.Transaction) bool {
	if item.Buffer > 0 {
		return false
	}
	updateBalanceOrder := balance.UnspentDeposits[0].Buffer <= item.Buffer // if the withdrawl empties the buffer of the top unspent transaction, then update the balance order
	balance.Total -= uint64(-item.Buffer)                                  //update the total balance
	//loop through the unspent deposits / transactions until the full amount of points are able to be withdrawn.
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
