package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

type AuthorityNode struct {
	PrivateKey           *ecdsa.PrivateKey
	TransactionPool      *blockchain.TransactionPool // Verwende den TransactionPool
	*Node                                            // Vererbung von Node
	LastBlockTimestamp   int64                       // Zeitstempel des letzten Blocks
	BlockCreationTrigger chan struct{}               // Kanal zum Auslösen der Blockerstellung
	mutex                sync.Mutex                  // Mutex zur Synchronisierung der Transaktionsverarbeitung
}

// Erstellt einen neuen AuthorityNode
func NewAuthorityNode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) *AuthorityNode {
	node := NewNode(privateKey, "localhost:8080")

	authorityNode := &AuthorityNode{
		PrivateKey:           privateKey,
		TransactionPool:      blockchain.NewTransactionPool(),
		Node:                 node,
		LastBlockTimestamp:   time.Now().Unix(),
		BlockCreationTrigger: make(chan struct{}),
		mutex:                sync.Mutex{},
	}

	// Erstelle den Genesis-Block
	authorityNode.Blockchain = blockchain.NewBlockchain(privateKey)

	go authorityNode.StartBlockGenerator()

	return authorityNode
}

func (a *AuthorityNode) AddTransaction(transaction *blockchain.Transaction) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Füge die Transaktion zum TransactionPool hinzu
	if err := a.TransactionPool.AddTransactionToPool(transaction); err != nil {
		return fmt.Errorf("error adding transaction to pool: %v", err)
	}

	// Überprüfe, ob die Anzahl der Transaktionen im Pool den Schwellenwert für die Blockerstellung erreicht
	if len(a.TransactionPool.GetTransactionsFromPool()) >= 5 {
		select {
		case a.BlockCreationTrigger <- struct{}{}:
			fmt.Println("BlockCreationTrigger was signalled")
		default:
			fmt.Println("BlockCreationTrigger already sent")
		}
	}

	return nil
}

func (a *AuthorityNode) CreateBlock() (*blockchain.Block, error) {
	// Sperre den Zugriff auf den TransactionPool
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Hole alle Transaktionen aus dem Transaktionspool
	pendingTransactions := a.TransactionPool.GetTransactionsFromPool()
	if len(pendingTransactions) < 1 {
		return nil, fmt.Errorf("not enough transactions to create a new block")
	}

	// Validierung jeder Transaktion vor dem Hinzufügen zum Block
	for _, tx := range pendingTransactions {
		doctorPublicKey, err := utils.DeserializePublicKey(tx.Doctor)
		if err != nil {
			return nil, fmt.Errorf("couldn't deserialize doctor public key: %v", err)
		}

		err = tx.ValidateTransaction(doctorPublicKey)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction in pool: %v", err)
		}
	}

	// Erstelle einen neuen Block mit den Transaktionen aus dem Pool
	newBlock := &blockchain.Block{
		ID:           uint64(len(a.Blockchain.Blocks) + 1),
		PreviousHash: a.Blockchain.Blocks[len(a.Blockchain.Blocks)-1].Hash,
		Transactions: pendingTransactions,
		Timestamp:    time.Now().Unix(),
	}

	// Berechne den Hash und signiere den Block
	hash, err := newBlock.CalculateHash()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %v", err)
	}
	newBlock.Hash = hash

	newBlock.SignBlock(a.PrivateKey)

	// Füge den Block zur Blockchain hinzu
	if err := a.AddBlockToBlockchain(newBlock); err != nil {
		return nil, fmt.Errorf("failed to add block to blockchain: %v", err)
	}

	// Entferne die Transaktionen aus dem Pool nach erfolgreicher Blockerstellung
	for _, tx := range pendingTransactions {
		txHash := fmt.Sprintf("%x", tx.Hash)
		a.TransactionPool.RemoveTransactionFromPool(txHash)
	}

	for _, tx := range newBlock.Transactions {
		patientHash := base64.StdEncoding.EncodeToString(tx.Patient)

		if _, exists := a.Patients[patientHash]; !exists {
			a.Patients[patientHash] = PatientData{
				Transactions: make(map[string]*blockchain.Transaction),
			}
		}

		a.Patients[patientHash].Transactions[base64.StdEncoding.EncodeToString(tx.Hash)] = tx
		fmt.Println(tx.Patient)
	}

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
	if !utils.VerifySignature(&a.PrivateKey.PublicKey, block.Hash, block.Signature.R, block.Signature.S) {
		return fmt.Errorf("invalid signature for transaction hash %x", block.Hash)
	}

	return nil
}

// check if conditions are met every 5th minute with sleep
func (a *AuthorityNode) StartBlockGenerator() {
	for {
		select {
		case <-time.After(5 * time.Minute):
			// Alle 5 Minuten versuchen, einen neuen Block zu erstellen
			if _, err := a.CreateBlock(); err != nil {
				fmt.Printf("Error creating a block: %v\n", err)
			}
		case <-a.BlockCreationTrigger:
			// Wenn eine neue Transaktion hinzugefügt wurde, versuche, einen Block zu erstellen
			if _, err := a.CreateBlock(); err != nil {
				fmt.Printf("Error creating a block: %v\n", err)
			}
		}
	}
}
