package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	ID           uint64
	Hash         []byte
	PreviousHash []byte
	Transactions []*Transaction
	Timestamp    int64
	Signature    []byte
}

func NewBlock(transactions []*Transaction, previousHash []byte, id uint64) (*Block, error) {
	block := &Block{
		ID:           id,
		PreviousHash: previousHash,
		Transactions: transactions,
		Timestamp:    GetCurrentTimestamp(),
	}

	hash, err := block.CalculateHash()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate block hash: %v", err)
	}
	block.Hash = hash

	return block, nil
}

func (b *Block) CalculateHash() ([]byte, error) {
	blockBytes, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block: %v", err)
	}

	hash := sha256.Sum256(blockBytes)
	return hash[:], nil // Rückgabe als Slice []byte
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func (b *Block) SerializeBlock() ([]byte, error) {
	blockBytes, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block: %v", err)
	}
	return blockBytes, nil
}

func DeserializeBlock(data []byte) (*Block, error) {
	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, fmt.Errorf("failed to deserialize block: %v", err)
	}
	return &block, nil
}

func (b *Block) PrintBlock() {
	fmt.Printf("Block ID: %d\n", b.ID)
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Printf("Previous Hash: %x\n", b.PreviousHash)
	fmt.Printf("Timestamp: %d\n", b.Timestamp)
	fmt.Printf("Number of Transactions: %d\n", len(b.Transactions))
	fmt.Println("Transactions:")
	for _, tx := range b.Transactions {
		tx.PrintTransaction()
	}
}
