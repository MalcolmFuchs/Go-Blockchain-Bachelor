package cmd

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

type AuthorityNode struct {
	PrivateKey           ed25519.PrivateKey // Private Key des Authority Nodes zur Signierung
	PublicKey            ed25519.PublicKey
	PendingTransactions  []*blockchain.Transaction // Liste der ausstehenden Transaktionen
	*Node                                          // Referenz auf den allgemeinen Node
	LastBlockTimestamp   int64                     // Zeitstempel des zuletzt erstellten Blocks
	BlockCreationTrigger chan struct{}
	mutex                sync.Mutex
}

func NewAuthorityNode(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) *AuthorityNode {
	node := NewNode(privateKey, publicKey, "localhost:8080")

	authorityNode := &AuthorityNode{
		PrivateKey:           privateKey,
		PublicKey:            publicKey,
		PendingTransactions:  []*blockchain.Transaction{},
		Node:                 node,
		LastBlockTimestamp:   time.Now().Unix(),
		BlockCreationTrigger: make(chan struct{}, 1),
	}

	// Create Genesis block
	authorityNode.Blockchain = blockchain.NewBlockchain(privateKey)

	go authorityNode.StartBlockGenerator()

	return authorityNode
}

func (a *AuthorityNode) AddTransaction(transaction *blockchain.Transaction) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if len(a.PendingTransactions) >= 5 {
		select {
		case a.BlockCreationTrigger <- struct{}{}:
			fmt.Println("BlockCreationTrigger was signalised")
		default:
			fmt.Println("BlockCreationTrigger been sent")
		}
	}

	a.PendingTransactions = append(a.PendingTransactions, transaction)
	fmt.Printf("Transaction %x added to pending transactions\n", transaction.Hash)
}

func (a *AuthorityNode) CreateBlock() (*blockchain.Block, error) {
	if len(a.PendingTransactions) < 1 {
		return nil, fmt.Errorf("not enough transactions to create a new block")
	}

	// Erstelle einen neuen Block mit den ausstehenden Transaktionen
	newBlock := &blockchain.Block{
		ID:           uint64(len(a.Blockchain.Blocks) + 1),
		PreviousHash: a.Blockchain.Blocks[len(a.Blockchain.Blocks)-1].Hash,
		Transactions: a.PendingTransactions,
		Timestamp:    time.Now().Unix(),
	}

	// Berechne den Hash des Blocks
	hash, err := newBlock.CalculateHash()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %v", err)
	}
	newBlock.Hash = hash

	// Signiere den Block mit dem Private Key des Authority Nodes
	newBlock.Signature = ed25519.Sign(a.PrivateKey, newBlock.Hash)

	// Füge den Block zur Blockchain hinzu
	if err := a.AddBlockToBlockchain(newBlock); err != nil {
		return nil, fmt.Errorf("failed to add block to blockchain: %v", err)
	}

	// Leere die Liste der ausstehenden Transaktionen
	a.PendingTransactions = []*blockchain.Transaction{}

	return newBlock, nil
}

func (a *AuthorityNode) AddBlockToBlockchain(block *blockchain.Block) error {

	if err := a.ValidateBlock(block); err != nil {
		return fmt.Errorf("failed to validate block: %v", err)
	}

	a.Blockchain.Blocks = append(a.Blockchain.Blocks, block)
	hashString := fmt.Sprintf("%x", block.Hash)
	a.Blockchain.BlockMap[hashString] = block

	a.LastBlockTimestamp = block.Timestamp

	fmt.Printf("Block with ID %d added to the blockchain\n", block.ID)
	return nil
}

func (a *AuthorityNode) CheckAndCreateBlock() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if len(a.PendingTransactions) >= 5 || (time.Now().Unix()-a.LastBlockTimestamp >= 300 && len(a.PendingTransactions) > 0) {
		_, err := a.CreateBlock()
		if err != nil {
			return fmt.Errorf("failed to create block: %v", err)
		}
	}

	return nil
}

func (a *AuthorityNode) ValidateBlock(block *blockchain.Block) error {
	// Erstelle eine temporäre Kopie des Blocks ohne Hash und Signatur
	tempBlock := *block
	tempBlock.Hash = nil
	tempBlock.Signature = nil

	// Berechne den Hash aus der temporären Blockkopie
	blockBytes, err := json.Marshal(tempBlock)
	if err != nil {
		return fmt.Errorf("failed to serialize block: %v", err)
	}
	calculatedHash := sha256.Sum256(blockBytes)

	// Überprüfe, ob der berechnete Hash mit dem gespeicherten Hash übereinstimmt
	if !bytes.Equal(calculatedHash[:], block.Hash) {
		return fmt.Errorf("invalid block hash for block ID %d", block.ID)
	}

	// Überprüfe, ob die Signatur gültig ist
	if !ed25519.Verify(a.Node.TrustedPublicKey, block.Hash, block.Signature) {
		return fmt.Errorf("invalid signature for block ID %d", block.ID)
	}

	return nil
}

func (a *AuthorityNode) GetPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	publicKeyHex := hex.EncodeToString(a.PublicKey)
	response := map[string]string{
		"publicKey": publicKeyHex,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// check if conditions are met every 5th minute with sleep
func (a *AuthorityNode) StartBlockGenerator() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := a.CheckAndCreateBlock()
			if err != nil {
				fmt.Printf("Error creating a block: %v\n", err)
			}
		default:
			if len(a.PendingTransactions) >= 5 {
				err := a.CheckAndCreateBlock()
				if err != nil {
					fmt.Printf("Error creating a block: %v\n", err)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}
