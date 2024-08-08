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
				Type:     "", // Wenn Typen leer bleiben, ersetzen Sie dies nach Bedarf
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
		patientData, exists := bc.Patients[id]
		if !exists {
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

		newBlock := bc.createBlock(newPatientRecord)
		dataToSign := fmt.Sprintf("%d%s%s%s", newBlock.Index, newBlock.Timestamp, newBlock.PatientData.PersonalData.InsuranceNumber, newBlock.PrevHash)
		r, s := SignData(dataToSign, node.PrivateKey)
		newBlock.SignatureR = r
		newBlock.SignatureS = s

		bc.validateAndAddBlock(newBlock, node)
	}
}

func (bc *Blockchain) GetMedicalRecords(id string, passphrase string) ([]MedicalRecord, error) {
	var medicalRecords []MedicalRecord

	bc.Mu.Lock()
	defer bc.Mu.Unlock()

	for _, block := range bc.Blocks {
		if block.PatientData.PersonalData.ID == id {
			fmt.Printf("Found block with patient data: %+v\n", block.PatientData)

			for _, encRecord := range block.PatientData.MedicalRecords {
				decryptedNotes := Decrypt(encRecord.Notes, passphrase)
				var decryptedRecord EncryptedMedicalRecord
				err := json.Unmarshal([]byte(decryptedNotes), &decryptedRecord)
				if err != nil {
					return nil, err
				}

				decryptedRecord.Date = Decrypt(decryptedRecord.Date, passphrase)
				decryptedRecord.Type = Decrypt(decryptedRecord.Type, passphrase)
				decryptedRecord.Provider = Decrypt(decryptedRecord.Provider, passphrase)
				decryptedRecord.Notes = Decrypt(decryptedRecord.Notes, passphrase)
				decryptedRecord.Results = Decrypt(decryptedRecord.Results, passphrase)

				recordDate, err := time.Parse(time.RFC3339, decryptedRecord.Date)
				if err != nil {
					return nil, err
				}

				medicalRecords = append(medicalRecords, MedicalRecord{
					Date:     recordDate,
					Type:     decryptedRecord.Type,
					Provider: decryptedRecord.Provider,
					Notes:    decryptedRecord.Notes,
					Results:  decryptedRecord.Results,
				})
			}
		}
	}

	if len(medicalRecords) == 0 {
		return nil, fmt.Errorf("no records found or access denied")
	}

	return medicalRecords, nil
}
