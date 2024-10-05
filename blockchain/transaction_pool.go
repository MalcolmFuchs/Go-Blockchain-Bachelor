package blockchain

import (
	"fmt"
)

type TransactionPool struct {
	Transactions map[string]*Transaction
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: make(map[string]*Transaction),
	}
}

func (tp *TransactionPool) AddTransactionToPool(transaction *Transaction) error {

	transactionHash := fmt.Sprintf("%x", transaction.Hash)
	if _, exists := tp.Transactions[transactionHash]; exists {
		return fmt.Errorf("transaction already exists in the pool")
	}

	tp.Transactions[transactionHash] = transaction
	fmt.Printf("Transaction %x added to the pool\n", transaction.Hash)

	return nil
}

func (tp *TransactionPool) RemoveTransactionFromPool(transactionHash string) error {
	if _, exists := tp.Transactions[transactionHash]; !exists {
		return fmt.Errorf("transaction %s does not exist in the pool", transactionHash)
	}

	delete(tp.Transactions, transactionHash)
	fmt.Printf("Transaction %s removed from the pool\n", transactionHash)

	return nil
}

func (tp *TransactionPool) GetTransactionsFromPool() []*Transaction {
	var transactions []*Transaction

	for _, transaction := range tp.Transactions {
		transactions = append(transactions, transaction)
	}

	return transactions
}
