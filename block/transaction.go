package blockchain

import "fmt"

func (bc *Blockchain) AddTransactionToPool(record EncryptedPatientRecord) {
	bc.PoolMu.Lock()
	defer bc.PoolMu.Unlock()

	bc.TransactionPool = append(bc.TransactionPool, record)
	fmt.Println("Encrypted transaction added to pool")

	if len(bc.TransactionPool) >= 10 {
		bc.CreateBlock()
	}
}

func (bc *Blockchain) AddEncryptedRecord(record EncryptedPatientRecord) {
	bc.PoolMu.Lock()
	defer bc.PoolMu.Unlock()

	bc.TransactionPool = append(bc.TransactionPool, record)
	fmt.Println("Encrypted transaction added to pool")
}

func (bc *Blockchain) ValidateTransaction(record EncryptedPatientRecord) bool {
	if record.PatientID == "" {
		fmt.Println("Invalid Transaction: Missing Patient ID")
		return false
	}

	fmt.Println("Transaction is valid.")
	return true
}

func (bc *Blockchain) ProcessTransaction(record EncryptedPatientRecord) {
	if bc.ValidateTransaction(record) {
		bc.AddTransactionToPool(record)
	}
}
