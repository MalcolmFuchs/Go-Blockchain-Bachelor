package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

func SendMedicalRecord(patientID string, record MedicalRecord, encryptionKey []byte, authorityNodeURL string) error {
	encryptedRecord := EncryptedMedicalRecord{}
	var err error

	encryptedRecord.Date, err = utils.Encrypt(record.Date, encryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting date: %v", err)
	}

	encryptedRecord.Type, err = utils.Encrypt(record.Type, encryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting type: %v", err)
	}

	encryptedRecord.Provider, err = utils.Encrypt(record.Provider, encryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting provider: %v", err)
	}

	encryptedRecord.Notes, err = utils.Encrypt(record.Notes, encryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting notes: %v", err)
	}

	encryptedRecord.Results, err = utils.Encrypt(record.Results, encryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting results: %v", err)
	}

	transaction := MedicalRecordTransaction{
		PatientID:       patientID,
		EncryptedRecord: encryptedRecord,
	}

	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("error serializing transaction: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/addTransaction", authorityNodeURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending transaction: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add medical record, status: %v", resp.Status)
	}

	fmt.Println("Encrypted medical record sent successfully.")
	return nil
}
