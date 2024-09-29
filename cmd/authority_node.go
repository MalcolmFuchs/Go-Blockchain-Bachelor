package cmd

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

type AuthorityNode struct {
	PrivateKey          ed25519.PrivateKey        // Private Key des Authority Nodes zur Signierung
	Blockchain          *blockchain.Blockchain    // Referenz auf die Blockchain-Struktur
	PendingTransactions []*blockchain.Transaction // Liste der ausstehenden Transaktionen
	Node                *Node                     // Referenz auf den allgemeinen Node
	LastBlockTimestamp  int64                     // Zeitstempel des zuletzt erstellten Blocks
}

func NewAuthorityNode(privateKey ed25519.PrivateKey, node *Node) *AuthorityNode {
	return &AuthorityNode{
		PrivateKey:          privateKey,
		Blockchain:          blockchain.NewBlockchain(privateKey),
		PendingTransactions: []*blockchain.Transaction{},
		Node:                node,
		LastBlockTimestamp:  time.Now().Unix(),
	}
}

func (a *AuthorityNode) AddTransaction(transaction *blockchain.Transaction) {
	a.PendingTransactions = append(a.PendingTransactions, transaction)
	fmt.Printf("Transaction %x added to pending transactions\n", transaction.Hash)
}

func (a *AuthorityNode) CreateBlock() (*blockchain.Block, error) {
	if len(a.PendingTransactions) < 10 {
		return nil, fmt.Errorf("not enough transactions to create a new block")
	}

	newBlock := &blockchain.Block{
		ID:           uint64(len(a.Blockchain.Blocks) + 1),
		PreviousHash: a.Blockchain.Blocks[len(a.Blockchain.Blocks)-1].Hash,
		Transactions: a.PendingTransactions,
		Timestamp:    time.Now().Unix(),
	}

	blockBytes, err := json.Marshal(newBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block: %v", err)
	}
	hash := sha256.Sum256(blockBytes)
	newBlock.Hash = hash[:]

	newBlock.Signature = ed25519.Sign(a.PrivateKey, newBlock.Hash)

	if err := a.AddBlockToBlockchain(newBlock); err != nil {
		return nil, fmt.Errorf("failed to add block to blockchain: %v", err)
	}

	a.PendingTransactions = []*blockchain.Transaction{}

	fmt.Printf("New block created with ID %d and hash %x\n", newBlock.ID, newBlock.Hash)
	return newBlock, nil
}

func (a *AuthorityNode) AddBlockToBlockchain(block *blockchain.Block) error {
	a.Blockchain.Blocks = append(a.Blockchain.Blocks, block)
	hashString := fmt.Sprintf("%x", block.Hash)
	a.Blockchain.BlockMap[hashString] = block

	a.LastBlockTimestamp = block.Timestamp

	fmt.Printf("Block with ID %d added to the blockchain\n", block.ID)
	return nil
}

func (a *AuthorityNode) CheckAndCreateBlock() error {
	if len(a.PendingTransactions) >= 10 || time.Now().Unix()-a.LastBlockTimestamp >= 300 {
		_, err := a.CreateBlock()
		if err != nil {
			return fmt.Errorf("failed to create block: %v", err)
		}
	}

	return nil
}

func (a *AuthorityNode) ValidateBlock(block *blockchain.Block) error {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to serialize block: %v", err)
	}
	calculatedHash := sha256.Sum256(blockBytes)

	if fmt.Sprintf("%x", calculatedHash[:]) != fmt.Sprintf("%x", block.Hash) {
		return fmt.Errorf("invalid block hash for block ID %d", block.ID)
	}

	if !ed25519.Verify(a.Node.TrustedPublicKey, block.Hash, block.Signature) {
		return fmt.Errorf("invalid signature for block ID %d", block.ID)
	}

	return nil
}
