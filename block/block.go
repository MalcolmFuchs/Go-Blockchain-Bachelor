package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	PrevHash      string
	Transactions  []Transaction
	Timestamp     time.Time
	Hash          string
	AuthoritySign []byte
}

func NewBlock(transactions []Transaction, prevHash string) Block {
	block := Block{
		PrevHash:     prevHash,
		Transactions: transactions,
		Timestamp:    time.Now(),
	}

	block.Hash = calculateHash(block)

	return block
}

// Berechne den Hash eines Blocks (basierend auf den Inhalten des Blocks)
func calculateHash(block Block) string {
	record := block.PrevHash + block.Timestamp.String() + hashTransactions(block.Transactions)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func hashTransactions(transactions []Transaction) string {
	var txHashes string
	for _, tx := range transactions {
		txHashes += tx.ID // Annahme: jede Transaktion hat eine ID
	}
	return txHashes
}
