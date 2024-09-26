package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Transaction struct {
	ID                   string
	EncryptedPatientData []byte
	EncryptedAESKey      []byte
	PatientSign          []byte
	DoctorSign           []byte
	Timestamp            time.Time
}

func NewTransaction(encryptedData, encryptedAESKey, patientSign []byte) Transaction {
	tx := Transaction{
		ID:                   generateTransactionID(),
		EncryptedPatientData: encryptedData,
		EncryptedAESKey:      encryptedAESKey,
		PatientSign:          patientSign,
		Timestamp:            time.Now(),
	}
	return tx
}

func generateTransactionID() string {
	h := sha256.New()
	h.Write([]byte(time.Now().String()))
	return hex.EncodeToString(h.Sum(nil))
}
