package components

import (
	"encoding/json"
	"fmt"
	"time"
)

func (bc *Blockchain) AddMedicalRecord(id string, newRecord MedicalRecord, passphrase string) {
	var patientFound bool

	node, err := bc.GetNextNode()
	if err != nil {
		fmt.Println("Error getting next node:", err)
		return
	}
	fmt.Println("Block added by " + node.Name)

	encryptedDate := Encrypt(newRecord.Date.Format(time.RFC3339), passphrase)
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

		if block.PatientData.PersonalData.ID == id {
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
		patientData, exits := bc.Patients[id]
		if !exits {
			fmt.Println("Patient not found")
			return
		}

		encryptedBirthDate := Encrypt(patientData.BirthDate.Format(time.RFC3339), passphrase)
		encryptedPatientData := EncryptedPersonalData{
			FirstName:       Encrypt(patientData.FirstName, passphrase),
			LastName:        Encrypt(patientData.LastName, passphrase),
			BirthDate:       encryptedBirthDate,
			InsuranceNumber: Encrypt(patientData.InsuranceNumber, passphrase),
		}

		newPatientRecord := PatientRecord{
			PersonalData: PersonalData{
				FirstName:       encryptedPatientData.FirstName,
				LastName:        encryptedPatientData.LastName,
				BirthDate:       patientData.BirthDate,
				InsuranceNumber: encryptedPatientData.InsuranceNumber,
				ID:              patientData.ID,
			},
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
			Timestamp:   time.Now(),
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

func (bc *Blockchain) GetMedicalRecords(id string, passphrase string) []MedicalRecord {

	for _, block := range bc.Blocks {
		if block.PatientData.PersonalData.ID == id {

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

				decryptedDate := Decrypt(decryptedRecord.Date, passphrase)
				date, err := time.Parse(time.RFC3339, decryptedDate)
				if err != nil {
					fmt.Println("Error parsing decrypted date:", err)
					continue
				}

				decryptedRecords = append(decryptedRecords, MedicalRecord{
					Date:     date,
					Type:     Decrypt(decryptedRecord.Type, passphrase),
					Provider: Decrypt(decryptedRecord.Provider, passphrase),
					Notes:    Decrypt(decryptedRecord.Notes, passphrase),
					Results:  Decrypt(decryptedRecord.Results, passphrase),
				})
			}
			return decryptedRecords
		}
	}
	return nil
}
