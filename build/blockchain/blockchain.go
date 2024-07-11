package blockchain

import (
	"time"
)

func CreateBlockchain(authorityNodes []AuthorityNode) Blockchain {
	genesisBlock := Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}

	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		authorityNodes,
	}
}

func (t *Transaction) validateTransaction() bool {
	if t.PatientID == "" {
		return false
	}

	record := t.Record
	if record.PersonalData.FirstName == "" || record.PersonalData.LastName == "" ||
		record.PersonalData.BirthDate.After(time.Now()) || len(record.PersonalData.InsuranceNumber) != 10 {
		return false
	}

	return true
}

func (b Blockchain) IsValid() bool {
	for i := range b.Chain[1:] {
		prevBlock := b.Chain[i]
		currBlock := b.Chain[i+1]

		if currBlock.Hash != currBlock.BlockHash() || currBlock.PrevHash != prevBlock.Hash {
			return false
		}

		// Überprüfen Sie die Transaktionen in jedem Block
		for _, transaction := range currBlock.Transactions {
			if !transaction.validateTransaction() {
				return false
			}
		}
	}
	return true
}

func (b *Blockchain) AddBlock(block Block) {
	b.Chain = append(b.Chain, block)
}
