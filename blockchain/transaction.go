package blockchain

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Transaction struct {
	Hash      []byte `json:"hash"`      // Der berechnete Hash der Transaktion
	Type      string `json:"type"`      // Typ der Transaktion (z.B. Checkup, Surgery)
	Notes     string `json:"notes"`     // Notizen zur Transaktion
	Results   string `json:"results"`   // Ergebnisse (z.B. Testergebnisse)
	Doctor    []byte `json:"doctor"`    // Public Key des Arztes als []byte
	Patient   []byte `json:"patient"`   // Public Key des Patienten als []byte
	Signature []byte `json:"signature"` // Signatur der Transaktion
	Key       []byte `json:"key"`       // AES-Schlüssel, verschlüsselt mit dem Public Key des Patienten
}

func NewTransaction(txType, notes, results, doctorPublicKeyHex, patientPublicKeyHex string, key []byte) (*Transaction, error) {
	// Wandelt Doctor- und Patient-Public-Key von Hex-String zu []byte um
	doctor, err := hex.DecodeString(doctorPublicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode doctor public key: %v", err)
	}

	patient, err := hex.DecodeString(patientPublicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode patient public key: %v", err)
	}

	tx := &Transaction{
		Type:    txType,
		Notes:   notes,
		Results: results,
		Doctor:  doctor,
		Patient: patient,
		Key:     key,
	}

	// Berechne den Hash der Transaktion
	hash, err := tx.CalculateHash()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate transaction hash: %v", err)
	}
	tx.Hash = hash

	return tx, nil
}

func (t *Transaction) CalculateHash() ([]byte, error) {
	// Serialisiere die Transaktion zu JSON, um den Hash daraus zu berechnen
	transactionBytes, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	hash := sha256.Sum256(transactionBytes)
	return hash[:], nil // Rückgabe des Hashes als Slice []byte
}

func (t *Transaction) SetHash() error {
	hash, err := t.CalculateHash()
	if err != nil {
		return fmt.Errorf("failed to calculate hash: %v", err)
	}
	t.Hash = hash // Setze den berechneten Hash in der Transaktion
	return nil
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
	// Berechne den Hash der Transaktion erneut
	hash, err := tx.CalculateHash()
	if err != nil {
		return fmt.Errorf("failed to calculate transaction hash: %v", err)
	}

	// Überprüfe, ob der Hash der Transaktion leer ist
	if len(hash) == 0 {
		return fmt.Errorf("transaction hash is empty")
	}

	// Überprüfe, ob der berechnete Hash mit dem gespeicherten Hash der Transaktion übereinstimmt
	if !bytes.Equal(hash, tx.Hash) {
		return fmt.Errorf("hash mismatch: calculated %x, stored %x", hash, tx.Hash)
	}

	// Verifiziere die Signatur mit dem Public Key
	if !ed25519.Verify(publicKey, tx.Hash, tx.Signature) {
		return fmt.Errorf("invalid signature for transaction hash %x", tx.Hash)
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