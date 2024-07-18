package components

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

func (b *Block) calculateHash() string {
	timestampString := b.Timestamp.Format(time.RFC3339)
	record := strconv.Itoa(b.Index) + timestampString + b.PatientData.PersonalData.InsuranceNumber + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func (bc *Blockchain) CreateGenesisBlock() {
	genesisBlock := Block{
		Index:       0,
		Timestamp:   time.Now(),
		PatientData: PatientRecord{},
		PrevHash:    "",
		Hash:        "",
	}
	genesisBlock.Hash = genesisBlock.calculateHash()

	bc.Blocks = append(bc.Blocks, genesisBlock)
}

func (bc *Blockchain) addBlock(newBlock Block) {
	bc.Blocks = append(bc.Blocks, newBlock)
}
