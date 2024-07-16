package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func (b *Block) calculateHash() string {
	record := string(b.Index) + b.Timestamp + b.PatientData.PersonalData.InsuranceNumber + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func createGenesisBlock() Block {
	genesisBlock := Block{
		Index:       0,
		Timestamp:   time.Now().String(),
		PatientData: PatientRecord{},
		PrevHash:    "",
		Hash:        "",
	}
	genesisBlock.Hash = genesisBlock.calculateHash()

	return genesisBlock
}
