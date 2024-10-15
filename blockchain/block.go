package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

type Block struct {
	ID           uint64
	Hash         []byte
	PreviousHash []byte
	Transactions []*Transaction
	Timestamp    int64
	Signature    *Signature `json:"signature"`
}

func (b *Block) CalculateHash() ([]byte, error) {
	// Erstelle eine temporäre Kopie des Blocks ohne Hash und Signatur
	tempBlock := *b
	tempBlock.Hash = nil
	tempBlock.Signature = nil

	blockBytes, err := json.Marshal(tempBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block: %v", err)
	}

	// Berechne den Hash aus den Blockdaten
	hash := sha256.Sum256(blockBytes)
	return hash[:], nil
}

func (b *Block) SignBlock(privateKey *ecdsa.PrivateKey) error {
	// Signiere den bereits berechneten Hash der Transaktion
	r, s, err := utils.SignTransaction(privateKey, b.Hash)
	if err != nil {
		return fmt.Errorf("failed to generate transaction signature: %v", err)
	}

	b.Signature = &Signature{
		R: r,
		S: s,
	}

	fmt.Printf("%v", b.Signature)

	return nil
}
