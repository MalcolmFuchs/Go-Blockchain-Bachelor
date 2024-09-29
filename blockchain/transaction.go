package blockchain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

type Transaction struct {
	Hash      []byte `json:"hash"`
	Data      []byte `json:"data"`
	Doctor    []byte `json:"doctor"`
	Patient   []byte `json:"patient"`
	Signature []byte `json:"signature"`
	Key       []byte `json:"key"`
	Timestamp int64  `json:"timestamp"`
}

type TransactionPool struct {
	Transactions map[string]*Transaction
	Mutex        sync.Mutex
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: make(map[string]*Transaction),
	}
}

func (tp *TransactionPool) AddTransactionToPool(transaction *Transaction) error {
	tp.Mutex.Lock()
	defer tp.Mutex.Unlock()

	transactionHashString := fmt.Sprintf("%x", transaction.Hash)

	if _, exists := tp.Transactions[transactionHashString]; exists {
		return fmt.Errorf("transaction with hash %s already exists in the pool", transactionHashString)
	}

	tp.Transactions[transactionHashString] = transaction
	fmt.Printf("Transaction added to pool with hash: %s\n", transactionHashString)
	return nil
}

func (tp *TransactionPool) RemoveTransactionFromPool(transactionHash string) error {
	tp.Mutex.Lock()
	defer tp.Mutex.Unlock()

	if _, exists := tp.Transactions[transactionHash]; !exists {
		return fmt.Errorf("transaction with hash %s does not exist in the pool", transactionHash)
	}

	delete(tp.Transactions, transactionHash)
	fmt.Printf("Transaction removed from pool with hash: %s\n", transactionHash)
	return nil
}

func (tp *TransactionPool) GetAllTransactions() []*Transaction {
	tp.Mutex.Lock()
	defer tp.Mutex.Unlock()

	transactions := make([]*Transaction, 0, len(tp.Transactions))
	for _, transaction := range tp.Transactions {
		transactions = append(transactions, transaction)
	}

	return transactions
}

func CreateTransaction(doctorPrivateKey ed25519.PrivateKey, doctorPublicKey, patientPublicKey, data []byte) (*Transaction, error) {

	aesKey := []byte("AESKey")
	encryptedData := data
	encryptedKey := aesKey

	transaction := &Transaction{
		Data:      encryptedData,
		Doctor:    doctorPublicKey,
		Patient:   patientPublicKey,
		Key:       encryptedKey,
		Timestamp: time.Now().Unix(),
	}

	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	hash := sha256.Sum256(transactionBytes)
	transaction.Hash = hash[:]
	transaction.Signature = ed25519.Sign(doctorPrivateKey, transaction.Hash)

	return transaction, nil
}

func SignTransaction(transaction *Transaction, privateKey ed25519.PrivateKey) error {
	if transaction == nil {
		return fmt.Errorf("transaction is nil")
	}

	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to serialize transaction: %v", err)
	}

	transactionHash := sha256.Sum256(transactionBytes)
	transaction.Hash = transactionHash[:] // Konvertiere [32]byte in []byte
	transaction.Signature = ed25519.Sign(privateKey, transaction.Hash)

	return nil
}

func VerifyTransaction(transaction *Transaction, publicKey ed25519.PublicKey) bool {
	if transaction == nil {
		return false
	}
	return ed25519.Verify(publicKey, transaction.Hash, transaction.Signature)
}

func CreateEncryptedTransaction(doctorPrivateKey ed25519.PrivateKey, doctorPublicKey, patientPublicKey, data []byte) (*Transaction, error) {

	aesKey := make([]byte, 32)
	_, err := utils.GenerateRandomBytes(len(aesKey))
	if err != nil {
		fmt.Println("Failed to generate AES key:", err)
	}

	encryptedData, _, err := utils.EncryptData(data, aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt patient data: %v", err)
	}

	encryptedAESKey, _, err := utils.EncryptWithPublicKey(aesKey, patientPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt AES key with patient's public key: %v", err)
	}

	transaction := &Transaction{
		Data:      encryptedData,
		Doctor:    doctorPublicKey,
		Patient:   patientPublicKey,
		Key:       encryptedAESKey,
		Timestamp: time.Now().Unix(),
	}

	if err := SignTransaction(transaction, doctorPrivateKey); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	return transaction, nil
}
