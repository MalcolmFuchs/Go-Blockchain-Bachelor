package blockchain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	ID           uint64         `json:"id"`
	Hash         []byte         `json:"hash"`
	PreviousHash []byte         `json:"previous_hash"`
	Transactions []*Transaction `json:"transactions"`
	Signature    []byte         `json:"signature"`
	Timestamp    int64          `json:"timestamp"`
}

type Blockchain struct {
	Blocks   []*Block          `json:"blocks"`
	BlockMap map[string]*Block `json:"block_map"`
}

func (bc *Blockchain) AddBlock(newBlock *Block, authorityPublicKey ed25519.PublicKey) error {
	// Validiere Block-Signatur mit Public Key des Authority Nodes
	if !ed25519.Verify(authorityPublicKey, newBlock.Hash, newBlock.Signature) {
		return fmt.Errorf("invalid block signature")
	}

	bc.Blocks = append(bc.Blocks, newBlock)

	// Block in BlockMap hinzufügen. Hash als String-Schlüssel.
	hashString := hex.EncodeToString(newBlock.Hash)
	bc.BlockMap[hashString] = newBlock

	return nil
}

func (bc *Blockchain) GetBlockchainData() ([]byte, error) {
	blockchainData, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize blockchain data: %v", err)
	}
	return blockchainData, nil
}

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

	blockHash := sha256.Sum256(blockBytes)
	genesisBlock.Hash = blockHash[:]
	genesisBlock.Signature = ed25519.Sign(authorityPrivateKey, genesisBlock.Hash)

	fmt.Println("Genesis Block created with ID 0 and hash:", genesisBlock.Hash)
	return genesisBlock, nil
}

func NewBlockchain(authorityPrivateKey ed25519.PrivateKey) (*Blockchain, error) {
	blockchain := &Blockchain{
		Blocks:   []*Block{},
		BlockMap: make(map[string]*Block),
	}

	genesisBlock, err := CreateGenesisBlock(authorityPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create genesis block: %v", err)
	}

	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)
	hashString := fmt.Sprintf("%x", genesisBlock.Hash)
	blockchain.BlockMap[hashString] = genesisBlock

	return blockchain, nil
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

	for i, block := range bc.Blocks {
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

func (bc *Blockchain) PrintBlockchain() {
	fmt.Println("Current Blockchain Overview:")

	// 1. Iteriere durch die Blöcke in der Blockchain und drucke Informationen aus
	for _, block := range bc.Blocks {
		fmt.Printf("Block ID: %d\n", block.ID)
		fmt.Printf("Hash: %x\n", block.Hash)
		if block.PreviousHash != nil {
			fmt.Printf("Previous Hash: %x\n", block.PreviousHash)
		} else {
			fmt.Println("Previous Hash: None (Genesis Block)")
		}
		fmt.Println("Transactions:")

		// 2. Iteriere durch die Transaktionen im Block
		for _, tx := range block.Transactions {
			fmt.Printf("  Transaction Hash: %x\n", tx.Hash)
			fmt.Printf("  Doctor: %x\n", tx.Doctor)
			fmt.Printf("  Patient: %x\n", tx.Patient)
			fmt.Printf("  Timestamp: %d\n", tx.Timestamp)
		}
		fmt.Println("-----------------------------------")
	}
}
