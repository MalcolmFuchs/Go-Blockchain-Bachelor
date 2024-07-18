package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"testing"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

// GenerateKeyPair generates a new public and private key pair
func GenerateKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return privateKey, privateKey.PublicKey
}

func TestGetMedicalRecords(t *testing.T) {
	// Initialize a new blockchain
	bc := blockchain.CreateBlockchain()

	// Generate a test passphrase
	passphrase := "testpassphrase"

	// Create test data
	patientID := "testpatient"
	medicalRecord := blockchain.MedicalRecord{
		Date:     time.Now(),
		Type:     "TestType",
		Provider: "TestProvider",
		Notes:    "TestNotes",
		Results:  "TestResults",
	}
	encryptedMedicalRecord := blockchain.EncryptedMedicalRecord{
		Date:     blockchain.Encrypt(medicalRecord.Date.Format(time.RFC3339), passphrase),
		Type:     blockchain.Encrypt(medicalRecord.Type, passphrase),
		Provider: blockchain.Encrypt(medicalRecord.Provider, passphrase),
		Notes:    blockchain.Encrypt(medicalRecord.Notes, passphrase),
		Results:  blockchain.Encrypt(medicalRecord.Results, passphrase),
	}

	// Create a test block
	block := blockchain.Block{
		Index:     1,
		Timestamp: time.Now(),
		PatientData: blockchain.PatientRecord{
			PersonalData: blockchain.PersonalData{
				ID: patientID,
			},
			MedicalRecords: []blockchain.MedicalRecord{
				medicalRecord,
			},
		},
		Hash:     "",
		PrevHash: "",
	}

	// Encrypt medical records and add them to the block
	var encryptedRecords []blockchain.MedicalRecord
	recordJSON, _ := json.Marshal(encryptedMedicalRecord)
	encryptedNotes := blockchain.Encrypt(string(recordJSON), passphrase)
	encryptedRecord := blockchain.MedicalRecord{
		Date:     medicalRecord.Date,
		Type:     medicalRecord.Type,
		Provider: medicalRecord.Provider,
		Notes:    encryptedNotes,
		Results:  medicalRecord.Results,
	}
	encryptedRecords = append(encryptedRecords, encryptedRecord)
	block.PatientData.MedicalRecords = encryptedRecords

	// Add the test block to the blockchain
	bc.Blocks = append(bc.Blocks, block)

	// Retrieve medical records using the function
	retrievedRecords := bc.GetMedicalRecords(patientID, passphrase)

	// Check if the retrieved records match the original ones
	if len(retrievedRecords) != 1 {
		t.Errorf("Expected 1 record, got %d", len(retrievedRecords))
	}

	retrievedRecord := retrievedRecords[0]
	if retrievedRecord.Type != medicalRecord.Type {
		t.Errorf("Expected Type %s, got %s", medicalRecord.Type, retrievedRecord.Type)
	}
	if retrievedRecord.Provider != medicalRecord.Provider {
		t.Errorf("Expected Provider %s, got %s", medicalRecord.Provider, retrievedRecord.Provider)
	}
	if retrievedRecord.Notes != medicalRecord.Notes {
		t.Errorf("Expected Notes %s, got %s", medicalRecord.Notes, retrievedRecord.Notes)
	}
	if retrievedRecord.Results != medicalRecord.Results {
		t.Errorf("Expected Results %s, got %s", medicalRecord.Results, retrievedRecord.Results)
	}
}
