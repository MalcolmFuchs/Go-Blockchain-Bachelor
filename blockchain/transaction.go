package blockchain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Transaction struct {
	Hash      []byte `json:"hash,omitempty"`
	Type      string `json:"type"`
	Notes     string `json:"notes"`
	Results   string `json:"results"`
	Doctor    string `json:"doctor"`
	Patient   string `json:"patient"`
	Signature []byte `json:"signature,omitempty"`
	Key       []byte `json:"key,omitempty"`
}

func NewTransaction(txType, notes, results, doctor, patient string, key []byte) (*Transaction, error) {
	tx := &Transaction{
		Type:    txType,
		Notes:   notes,
		Results: results,
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
	return hash[:], nil
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
	fmt.Printf("Doctor: %s\n", t.Doctor)
	fmt.Printf("Patient: %s\n", t.Patient)
	fmt.Printf("Type: %s\n", t.Type)
	fmt.Printf("Notes: %s\n", t.Notes)
	fmt.Printf("Results: %s\n", t.Results)
	fmt.Printf("Signature: %x\n", t.Signature)
}
