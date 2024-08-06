package components

import (
	"time"
)

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

func (bc *Blockchain) createBlock(transaction PatientRecord) Block {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		Index:       len(bc.Blocks),
		Timestamp:   time.Now(),
		PatientData: transaction,
		PrevHash:    prevBlock.Hash,
	}
	newBlock.Hash = bc.calculatBcHash(newBlock)
	return newBlock
}
