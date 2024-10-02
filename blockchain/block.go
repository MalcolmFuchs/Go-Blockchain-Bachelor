package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Block struct {
	ID           uint64
	Hash         []byte
	PreviousHash []byte
	Transactions []*Transaction
	Timestamp    int64
	Signature    []byte
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
