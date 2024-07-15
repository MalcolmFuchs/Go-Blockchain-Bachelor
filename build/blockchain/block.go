package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func (block Block) BlockHash() string {
	transactionsData, _ := json.Marshal(block.Transactions)
	blockData := block.PrevHash + string(transactionsData) + block.Timestamp.String() + strconv.Itoa(block.Index)
	blockHash := sha256.Sum256([]byte(blockData))

	return fmt.Sprintf("%x", blockHash)
}

func (node AuthorityNode) ValidateBlock(newBlock *Block, blockchain *Blockchain) bool {
	// Ist der Hash des neuen Blocks korrekt?
	if newBlock.Hash != newBlock.BlockHash() {
		return false
	}

	// Stimmt der vorherige Hash mit dem Hash des letzten Blocks in der Blockchain überein?
	lastBlock := blockchain.Chain[len(blockchain.Chain)-1]
	if newBlock.PrevHash != lastBlock.Hash {
		return false
	}

	return true
}

func NewBlock(transactions []Transaction, prevBlock *Block, allNodes []AuthorityNode, blockchain *Blockchain) *Block {
	b := new(Block)

	b.Index = prevBlock.Index + 1
	b.Timestamp = time.Now()
	b.Transactions = transactions
	b.PrevHash = prevBlock.Hash
	b.Hash = b.BlockHash()

	return b
}

func (node AuthorityNode) CreateBlock(transactions []Transaction, prevBlock *Block, allNodes []AuthorityNode, blockchain *Blockchain) *Block {
	b := NewBlock(transactions, prevBlock, allNodes, blockchain)

	// Senden des neuen Blocks an alle anderen Authority Nodes zur Validierung
	for _, otherNode := range allNodes {
		if otherNode.Id != node.Id && !otherNode.ValidateBlock(b, blockchain) {

			return nil
		}
	}

	return b
}
