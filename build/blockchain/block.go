package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	index     int
	timestamp time.Time
	hash      string
	data      map[string]interface{}
	prevHash  string
}

func (block Block) Hash() string {
	data, _ := json.Marshal(block.data)
	blockData := block.prevHash + string(data) + block.timestamp.String() + strconv.Itoa(block.index)
	blockHash := sha256.Sum256([]byte(blockData))

	return fmt.Sprintf("%x", blockHash)
}

func (node *AuthorityNode) createBlock(data map[string]interface{}, prevBlock Block) Block {
	var newBlock Block

	newBlock.index = prevBlock.index + 1
	newBlock.timestamp = time.Now()
	newBlock.data = data
	newBlock.prevHash = prevBlock.hash
	newBlock.hash = newBlock.Hash()

	return newBlock
}
