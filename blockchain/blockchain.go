package blockchain

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (bc *Blockchain) addBlock(newBlock Block) {
	bc.Blocks = append(bc.Blocks, newBlock)
}

func CreateBlockchain() *Blockchain {
	return &Blockchain{[]Block{createGenesisBlock()}, nil}
}

func (bc *Blockchain) validateAndAddBlock(newBlock Block, node *AuthorityNode) {
	if len(bc.Blocks) > 0 {
		lastBlock := bc.Blocks[len(bc.Blocks)-1]
		if newBlock.PrevHash != lastBlock.Hash {
			fmt.Println("Invalid block: PrevHash does not match")
			return
		}
	}

	newBlock.Hash = newBlock.calculateHash()
	dataToSign := fmt.Sprintf("%d%s%s%s", newBlock.Index, newBlock.Timestamp, newBlock.PatientData.PersonalData.InsuranceNumber, newBlock.PrevHash)
	r, s := SignData(dataToSign, node.PrivateKey)
	newBlock.SignatureR = r
	newBlock.SignatureS = s

	if VeryfiySignature(dataToSign, newBlock.SignatureR, newBlock.SignatureS, node.PublicKey) {
		bc.addBlock(newBlock)
		fmt.Println("Block added by", node.Name)
	} else {
		fmt.Println("Invalid signature. Block not added.")
	}
}

func (bc *Blockchain) AddMedicalRecord(insuranceNumber string, newRecord MedicalRecord, passphrase string) {
	var patientFound bool

	node, err := bc.GetRandomNode()
	if err != nil {
		fmt.Println("Error getting random node:", err)
		return
	}
	encryptedRecord := Encrypt(fmt.Sprintf("%v", newRecord), passphrase)

	for i, block := range bc.Blocks {
		if block.PatientData.PersonalData.InsuranceNumber == insuranceNumber {
			bc.Blocks[i].PatientData.MedicalRecords = append(bc.Blocks[i].PatientData.MedicalRecords, MedicalRecord{
				Date:     newRecord.Date,
				Type:     newRecord.Type,
				Provider: newRecord.Provider,
				Notes:    encryptedRecord,
				Results:  newRecord.Results,
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
		newPatientData := PersonalData{InsuranceNumber: insuranceNumber}
		newPatientRecord := PatientRecord{
			PersonalData: newPatientData,
			MedicalRecords: []MedicalRecord{{
				Date:     newRecord.Date,
				Type:     newRecord.Type,
				Provider: newRecord.Provider,
				Notes:    encryptedRecord,
				Results:  newRecord.Results,
			}},
		}
		newBlock := Block{
			Index:       len(bc.Blocks),
			Timestamp:   time.Now().String(),
			PatientData: newPatientRecord,
			PrevHash:    bc.Blocks[len(bc.Blocks)-1].Hash,
			Hash:        "",
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
			for i, record := range block.PatientData.MedicalRecords {
				decypredNotes := Decrypt(record.Notes, passphrase)
				block.PatientData.MedicalRecords[i].Notes = decypredNotes
			}
			return block.PatientData.MedicalRecords
		}
	}
	return nil
}

func (bc *Blockchain) GetRandomNode() (*AuthorityNode, error) {
	if len(bc.Nodes) == 0 {
		return nil, errors.New("no nodes available")
	}
	index := rand.Intn(len(bc.Nodes))
	return &bc.Nodes[index], nil
}
