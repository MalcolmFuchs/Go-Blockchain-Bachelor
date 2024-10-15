package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

type Transaction struct {
	Hash          []byte              `json:"hash"`
	EncryptedData utils.EncryptedData `json:"encryptedData"`
	Doctor        []byte              `json:"doctor"`
	Patient       []byte              `json:"patient"`
	Signature     []byte              `json:"signature"`
}

type TransactionData struct {
	Type    string `json:"type"`
	Notes   string `json:"notes"`
	Results string `json:"results"`
}

func PrepareTransactionData(txType, notes, results string) ([]byte, error) {
	data := TransactionData{
		Type:    txType,
		Notes:   notes,
		Results: results,
	}
	return json.Marshal(data)
}

func NewTransaction(txType, notes, results string, senderPrivKey *ecdsa.PrivateKey, recipientPubKey *ecdsa.PublicKey) (*Transaction, error) {
	// Bereite die Transaktionsdaten vor
	plaintext, err := PrepareTransactionData(txType, notes, results)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transaction data: %v", err)
	}

	senderEcdhPrivKey, err := utils.EcdsaPrivToEcdh(senderPrivKey)
	if err != nil {
		return nil, fmt.Errorf("Error during conversion from ecdsa to ecdh private key", err)
	}

	recipientEcdhPubKey, err := utils.EcdsaPubToEcdh(recipientPubKey)
	if err != nil {
		return nil, fmt.Errorf("Error during conversion from ecdsa to ecdh public key", err)
	}

	// Verschlüssele die Daten mit AES-GCM
	ciphertext, nonce, err := utils.EncryptData(senderEcdhPrivKey, recipientEcdhPubKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction data: %v", err)
	}

	tx := &Transaction{
		Doctor:  senderPrivKey.X.Bytes(),
		Patient: recipientPubKey.X.Bytes(),
		EncryptedData: utils.EncryptedData{
			Ciphertext: ciphertext,
			Nonce:      nonce,
		},
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
	// Erstelle eine temporäre Kopie der Transaktion ohne Hash und Signatur
	tempTx := *t
	tempTx.Hash = nil
	tempTx.Signature = nil

	transactionBytes, err := json.Marshal(tempTx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	hash := sha256.Sum256(transactionBytes)
	return hash[:], nil // Rückgabe des Hashes als Slice []byte
}

func SignTransaction(tx *Transaction, privateKey ed25519.PrivateKey) ([]byte, error) {
	// Signiere den bereits berechneten Hash der Transaktion
	signature := ed25519.Sign(privateKey, tx.Hash)
	tx.Signature = signature
	return signature, nil
}

func ValidateTransaction(tx *Transaction, publicKey ed25519.PublicKey) error {
	// Berechne den Hash der Transaktion erneut
	hash, err := tx.CalculateHash()
	if err != nil {
		return fmt.Errorf("failed to calculate transaction hash: %v", err)
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
