package blockchain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Transaction struct {
	Hash      []byte
	Data      []byte
	Doctor    []byte
	Patient   []byte
	Signature []byte
	Key       []byte // AES-Schlüssel, verschlüsselt mit dem Public Key des Patienten
}

func NewTransaction(data, doctor, patient, key []byte) (*Transaction, error) {
	tx := &Transaction{
		Data:    data,
		Doctor:  doctor,
		Patient: patient,
		Key:     key,
	}

	hash, err := tx.CalculateHash()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate transaction hash: %v", err)
	}
	tx.Hash = hash

	return tx, nil
}

func (t *Transaction) CalculateHash() ([]byte, error) {
	transactionBytes, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	hash := sha256.Sum256(transactionBytes)
	return hash[:], nil // Rückgabe als Slice []byte
}

func SignTransaction(tx *Transaction, privateKey []byte) ([]byte, error) {
	transactionBytes, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	hash := sha256.Sum256(transactionBytes)
	signature := ed25519.Sign(privateKey, hash[:])
	tx.Signature = signature

	return signature, nil
}

func ValidateTransaction(tx *Transaction, publicKey []byte) error {
	hash, err := tx.CalculateHash()
	if err != nil {
		return fmt.Errorf("failed to calculate transaction hash: %v", err)
	}

	if !ed25519.Verify(publicKey, hash, tx.Signature) {
		return fmt.Errorf("invalid signature for transaction")
	}

	return nil
}

func (t *Transaction) SerializeTransaction() ([]byte, error) {
	transactionBytes, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}
	return transactionBytes, nil
}

func DeserializeTransaction(data []byte) (*Transaction, error) {
	var tx Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}
	return &tx, nil
}

func (t *Transaction) PrintTransaction() {
	fmt.Printf("Transaction Hash: %x\n", t.Hash)
	fmt.Printf("Doctor: %x\n", t.Doctor)
	fmt.Printf("Patient: %x\n", t.Patient)
	fmt.Printf("Encrypted Data: %x\n", t.Data)
	fmt.Printf("Encrypted AES Key: %x\n", t.Key)
	fmt.Printf("Signature: %x\n", t.Signature)
}
