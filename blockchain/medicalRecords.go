package blockchain

import (
	"fmt"
	"time"
)

func (bc *Blockchain) AddMedicalRecord(insuranceNumber string, newRecord MedicalRecord, passphrase string) {
	var patientFound bool

	node, err := bc.GetRandomNode()
	if err != nil {
		fmt.Println("Error getting random node:", err)
		return
	}

	encryptedType := Encrypt(fmt.Sprintf("%v", newRecord.Type), passphrase)
	encryptedProvider := Encrypt(fmt.Sprintf("%v", newRecord.Provider), passphrase)
	encryptedNotes := Encrypt(fmt.Sprintf("%v", newRecord.Notes), passphrase)
	encryptedResults := Encrypt(fmt.Sprintf("%v", newRecord.Results), passphrase)

	for i, block := range bc.Blocks {

		if block.PatientData.PersonalData.InsuranceNumber == insuranceNumber {
			bc.Blocks[i].PatientData.MedicalRecords = append(bc.Blocks[i].PatientData.MedicalRecords, MedicalRecord{
				Date:     newRecord.Date,
				Type:     encryptedType,
				Provider: encryptedProvider,
				Notes:    encryptedNotes,
				Results:  encryptedResults,
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
				Type:     encryptedType,
				Provider: encryptedProvider,
				Notes:    encryptedNotes,
				Results:  encryptedResults,
			}},
		}
		newBlock := Block{
			Index:       len(bc.Blocks),
			Timestamp:   time.Now().String(),
			PatientData: newPatientRecord,
			PrevHash:    bc.Blocks[len(bc.Blocks)-1].Hash,
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
			fmt.Println(block.PatientData.MedicalRecords)
			for i, record := range block.PatientData.MedicalRecords {
				decypredType := Decrypt(record.Type, passphrase)
				block.PatientData.MedicalRecords[i].Type = decypredType
				decypredProvider := Decrypt(record.Provider, passphrase)
				block.PatientData.MedicalRecords[i].Provider = decypredProvider
				decypredNotes := Decrypt(record.Notes, passphrase)
				block.PatientData.MedicalRecords[i].Notes = decypredNotes
				decypredResults := Decrypt(record.Results, passphrase)
				block.PatientData.MedicalRecords[i].Results = decypredResults
			}
			return block.PatientData.MedicalRecords
		}
	}
	return nil
}
