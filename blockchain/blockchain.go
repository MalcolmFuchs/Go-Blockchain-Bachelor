package blockchain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Blockchain represents the structure of the blockchain containing all blocks and a map for quick lookup
type Blockchain struct {
	Blocks   []*Block          // Liste aller Blöcke in der Blockchain
	BlockMap map[string]*Block // Mapping von Block-Hash zu Block, um schnellen Zugriff zu ermöglichen
}

// NewBlockchain creates a new blockchain with a genesis block
func NewBlockchain(privateKey ed25519.PrivateKey) *Blockchain {
	// Erstelle den Genesis-Block und initialisiere die Blockchain
	genesisBlock, err := CreateGenesisBlock(privateKey)
	if err != nil {
		fmt.Printf("Failed to create genesis block: %v\n", err)
		return nil
	}

	blockchain := &Blockchain{
		Blocks:   []*Block{genesisBlock},
		BlockMap: map[string]*Block{hex.EncodeToString(genesisBlock.Hash): genesisBlock},
	}

	return blockchain
}

// CreateGenesisBlock creates the initial block of the blockchain
func CreateGenesisBlock(authorityPrivateKey ed25519.PrivateKey) (*Block, error) {
	genesisBlock := &Block{
		ID:           0,
		PreviousHash: nil,
		Transactions: []*Transaction{},
		Timestamp:    time.Now().Unix(),
	}

	blockBytes, err := json.Marshal(genesisBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize genesis block: %v", err)
	}

	hash := sha256.Sum256(blockBytes)
	genesisBlock.Hash = hash[:]

	// 4. Signiere den Genesis-Block mit dem Private Key des Authority Nodes
	genesisBlock.Signature = ed25519.Sign(authorityPrivateKey, genesisBlock.Hash)

	fmt.Println("Genesis Block created with ID 0 and hash:", genesisBlock.Hash)
	return genesisBlock, nil
}
