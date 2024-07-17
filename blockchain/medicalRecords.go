package blockchain

import (
	"encoding/json"
	"fmt"
	"time"
)

func (bc *Blockchain) AddMedicalRecord(insuranceNumber string, newRecord MedicalRecord, passphrase string) {
	var patientFound bool

	node, err := bc.GetNextNode()
	if err != nil {
		fmt.Println("Error getting next node:", err)
		return
	}

	fmt.Println("Block added by " + node.Name)

	encryptedDate := Encrypt(CustomDateToString(newRecord.Date), passphrase)

	encRecord := EncryptedMedicalRecord{
		Date:     encryptedDate,
		Type:     Encrypt(newRecord.Type, passphrase),
		Provider: Encrypt(newRecord.Provider, passphrase),
		Notes:    Encrypt(newRecord.Notes, passphrase),
		Results:  Encrypt(newRecord.Results, passphrase),
	}
	recordBytes, err := json.Marshal(encRecord)
	if err != nil {
		fmt.Println("Error marshaling medical record:", err)
		return
	}
	encryptedRecord := Encrypt(string(recordBytes), passphrase)

	for i, block := range bc.Blocks {

		if block.PatientData.PersonalData.InsuranceNumber == insuranceNumber {
			bc.Blocks[i].PatientData.MedicalRecords = append(bc.Blocks[i].PatientData.MedicalRecords, MedicalRecord{
				Date:     newRecord.Date,
				Type:     "",
				Provider: "",
				Notes:    encryptedRecord,
				Results:  "",
			})
			bc.Blocks[i].Hash = bc.Blocks[i].calculateHash()
			dataToSign := fmt.Sprintf("%d%s%s%s", bc.Blocks[i].Index, bc.Blocks[i].Timestamp, bc.Blocks[i].PatientData.PersonalData.InsuranceNumber, bc.Blocks[i].PrevHash)
			r, s := SignData(dataToSign, node.PrivateKey)
			bc.Blocks[i].SignatureR = r
			bc.Blocks[i].SignatureS = s
			patientFound = true
			break
		}
	}
	if !patientFound {
		patientData, exits := bc.Patients[insuranceNumber]
		if !exits {
			fmt.Println("Patient not found")
			return
		}
		newPatientRecord := PatientRecord{
			PersonalData: patientData,
			MedicalRecords: []MedicalRecord{{
				Date:     newRecord.Date,
				Type:     "",
				Provider: "",
				Notes:    encryptedRecord,
				Results:  "",
			}},
		}
		newBlock := Block{
			Index:       len(bc.Blocks),
			Timestamp:   time.Now().String(),
			PatientData: newPatientRecord,
			PrevHash:    "",
			Hash:        "",
		}
		if len(bc.Blocks) > 0 {
			newBlock.PrevHash = bc.Blocks[len(bc.Blocks)-1].Hash
		}
		newBlock.Hash = newBlock.calculateHash()
		dataToSign := fmt.Sprintf("%d%s%s%s", newBlock.Index, newBlock.Timestamp, newBlock.PatientData.PersonalData.InsuranceNumber, newBlock.PrevHash)
		r, s := SignData(dataToSign, node.PrivateKey)
		newBlock.SignatureR = r
		newBlock.SignatureS = s

		if VeryfiySignature(dataToSign, newBlock.SignatureR, newBlock.SignatureS, node.PublicKey) {
			bc.validateAndAddBlock(newBlock, node)
		}
	}
}

func (bc *Blockchain) GetMedicalRecords(insuranceNumber string, passphrase string) []MedicalRecord {
	for _, block := range bc.Blocks {
		if block.PatientData.PersonalData.InsuranceNumber == insuranceNumber {
			var decryptedRecords []MedicalRecord
			for _, record := range block.PatientData.MedicalRecords {
				decryptedData := Decrypt(record.Notes, passphrase)
				if decryptedData == "" {
					continue
				}
				var decryptedRecord EncryptedMedicalRecord
				err := json.Unmarshal([]byte(decryptedData), &decryptedRecord)
				if err != nil {
					fmt.Println("Error unmarshaling decrypted data:", err)
					continue
				}

				// Decrypt each field of the decryptedRecord
				decryptedDate := Decrypt(decryptedRecord.Date, passphrase)
				if decryptedDate == "" {
					continue
				}
				date, err := time.Parse(time.RFC3339, decryptedDate)
				if err != nil {
					fmt.Println("Error parsing decrypted date:", err)
					continue
				}

				decryptedType := Decrypt(decryptedRecord.Type, passphrase)
				decryptedProvider := Decrypt(decryptedRecord.Provider, passphrase)
				decryptedNotes := Decrypt(decryptedRecord.Notes, passphrase)
				decryptedResults := Decrypt(decryptedRecord.Results, passphrase)

				decryptedRecords = append(decryptedRecords, MedicalRecord{
					Date:     date,
					Type:     decryptedType,
					Provider: decryptedProvider,
					Notes:    decryptedNotes,
					Results:  decryptedResults,
				})
			}
			return decryptedRecords
		}
	}
	return nil
}
