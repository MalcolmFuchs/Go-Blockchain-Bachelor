package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

func (bc *Blockchain) CreateGenesisBlock() {
	genesisBlock := Block{
		Index:      0,
		Timestamp:  time.Now(),
		Hash:       "",
		PrevHash:   "",
		SignatureR: "",
		SignatureS: "",
	}

	genesisBlock.Hash = bc.CalculateHash(genesisBlock)

	r, s := utils.SignData(genesisBlock.Hash, bc.PrivateKey)
	genesisBlock.SignatureR = r
	genesisBlock.SignatureS = s

	bc.Blocks = append(bc.Blocks, genesisBlock)
	fmt.Println("Genesis Block created and signed.")
}

func (bc *Blockchain) CreateBlock() {
	bc.Mu.Lock()
	defer bc.Mu.Unlock()

	newBlock := Block{
		Index:       len(bc.Blocks),
		Timestamp:   time.Now(),
		PatientData: bc.TransactionPool,
		PrevHash:    bc.getPreviousHash(),
	}

	newBlock.Hash = bc.CalculateHash(newBlock)

	r, s := utils.SignData(newBlock.Hash, bc.PrivateKey)
	newBlock.SignatureR = r
	newBlock.SignatureS = s

	bc.Blocks = append(bc.Blocks, newBlock)
	fmt.Printf("Block %d created and signed\n", newBlock.Index)

	bc.TransactionPool = []EncryptedPatientRecord{}
}

func (bc *Blockchain) CalculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s", block.Index, block.Timestamp, block.PrevHash)
	hash := sha256.New()
	hash.Write([]byte(record))
	return hex.EncodeToString(hash.Sum(nil))
}

func (bc *Blockchain) getPreviousHash() string {
	if len(bc.Blocks) == 0 {
		return ""
	}
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	return lastBlock.Hash
}

func (bc *Blockchain) StartBlockTimer() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bc.Mu.Lock()
			if len(bc.TransactionPool) > 0 {
				bc.CreateBlock()
			}
			bc.Mu.Unlock()
		}
	}
}
