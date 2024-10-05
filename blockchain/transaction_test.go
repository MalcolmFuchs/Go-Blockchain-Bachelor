// transaction_test.go

package blockchain

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"
)

func generateRandomHexKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func TestValidateTransaction(t *testing.T) {
	// Schlüsselpaare für Arzt, Patient und Authority Node
	doctorPub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate doctor keys: %v", err)
	}

	patientPub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate patient keys: %v", err)
	}

	authorityPub, authorityPriv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate authority node keys: %v", err)
	}

	// Generiere einen zufälligen AES-Schlüssel und encodiere ihn als Hex
	validKeyHex, err := generateRandomHexKey(32) // 32 Bytes für AES-256
	if err != nil {
		t.Fatalf("Failed to generate random AES key: %v", err)
	}

	// Erstelle eine Transaktion (ohne Signatur)
	tx, err := NewTransaction(
		"Checkup",
		"Routine checkup",
		"All normal",
		hex.EncodeToString(doctorPub),
		hex.EncodeToString(patientPub),
		validKeyHex,
	)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Signiere die Transaktion mit dem Authority Node's Schlüssel
	_, err = SignTransaction(tx, authorityPriv)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	// Validierung der Transaktion sollte erfolgreich sein
	err = ValidateTransaction(tx, authorityPub)
	if err != nil {
		t.Errorf("Expected transaction to be valid, got error: %v", err)
	}

	// Ändere eine Transaktionseigenschaft, um die Validierung zu fehlschlagen
	tx.Notes = "Tampered notes"

	err = ValidateTransaction(tx, authorityPub)
	if err == nil {
		t.Errorf("Expected transaction to be invalid due to tampering, but validation passed")
	}
}
