package components

import (
	"time"
)

func (bc *Blockchain) AddTransactionToPool(transaction PatientRecord) {
	bc.poolMu.Lock()
	defer bc.poolMu.Unlock()
	bc.TransactionPool = append(bc.TransactionPool, transaction)
}

func (bc *Blockchain) ProcessTransactions() {
	for {
		bc.poolMu.Lock()
		if len(bc.TransactionPool) == 0 {
			bc.poolMu.Unlock()
			time.Sleep(time.Second + 10)
			continue
		}

		transaction := bc.TransactionPool[0]
		bc.TransactionPool = bc.TransactionPool[1:]
		bc.poolMu.Unlock()

		newBlock := bc.createBlock(transaction)
		bc.mu.Lock()
		bc.Blocks = append(bc.Blocks, newBlock)
		bc.mu.Unlock()
	}
}
