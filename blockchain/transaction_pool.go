package blockchain

import (
	"crypto/ed25519"
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

func (tp *TransactionPool) AddTransactionToPool(transaction *Transaction, doctorPublicKey ed25519.PublicKey) error {

	if !ed25519.Verify(doctorPublicKey, transaction.Hash, transaction.Signature) {
		fmt.Printf("Signature verification failed. Doctor's Public Key: %x, Hash: %x, Signature: %x\n", doctorPublicKey, transaction.Hash, transaction.Signature)
		return fmt.Errorf("invalid transaction signature for transaction hash %x", transaction.Hash)
	}

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

func (tp *TransactionPool) PrintPool() {
	fmt.Printf("Transaction Pool: %d transactions\n", len(tp.Transactions))
	for hash, transaction := range tp.Transactions {
		fmt.Printf("Transaction Hash: %s\n", hash)
		transaction.PrintTransaction()
	}
}
