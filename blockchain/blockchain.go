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

// AddBlock adds a new block to the blockchain and updates the BlockMap
func (bc *Blockchain) AddBlock(block *Block) error {
	// 1. Überprüfe, ob der Block mit der Blockchain konsistent ist
	if len(bc.Blocks) > 0 {
		lastBlock := bc.Blocks[len(bc.Blocks)-1]
		if hex.EncodeToString(block.PreviousHash) != hex.EncodeToString(lastBlock.Hash) {
			return fmt.Errorf("block's previous hash does not match the last block's hash")
		}
	}

	bc.Blocks = append(bc.Blocks, block)
	bc.BlockMap[hex.EncodeToString(block.Hash)] = block

	fmt.Printf("Block with ID %d and hash %x added to the blockchain\n", block.ID, block.Hash)
	return nil
}

// GetBlockchainData returns the serialized blockchain data as JSON
func (bc *Blockchain) GetBlockchainData() ([]byte, error) {
	blockchainData, err := json.Marshal(bc.Blocks)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize blockchain data: %v", err)
	}

	return blockchainData, nil
}

func (bc *Blockchain) ValidateBlock(block *Block, authorityPublicKey ed25519.PublicKey) error {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to serialize block: %v", err)
	}
	calculatedHash := sha256.Sum256(blockBytes)

	if hex.EncodeToString(calculatedHash[:]) != hex.EncodeToString(block.Hash) {
		return fmt.Errorf("invalid block hash for block ID %d", block.ID)
	}

	if !ed25519.Verify(authorityPublicKey, block.Hash, block.Signature) {
		return fmt.Errorf("invalid signature for block ID %d", block.ID)
	}

	return nil
}

func (bc *Blockchain) ValidateBlockchain(authorityPublicKey ed25519.PublicKey) error {
	// 1. Iteriere durch alle Blöcke der Blockchain
	for i, block := range bc.Blocks {
		// 2. Validiere den aktuellen Block
		if err := bc.ValidateBlock(block, authorityPublicKey); err != nil {
			return fmt.Errorf("block validation failed: %v", err)
		}

		// 3. Überprüfe, ob der Hash des vorherigen Blocks korrekt ist (außer beim Genesis-Block)
		if i > 0 && hex.EncodeToString(block.PreviousHash) != hex.EncodeToString(bc.Blocks[i-1].Hash) {
			return fmt.Errorf("previous hash mismatch at block ID %d", block.ID)
		}
	}

	fmt.Println("Blockchain validation successful. All blocks are valid.")
	return nil
}
