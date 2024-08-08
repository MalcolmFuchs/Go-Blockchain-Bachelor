package components

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"time"
)

func (bc *Blockchain) AddTransactionToPool(transaction PatientRecord, r, s *big.Int, publicKey ecdsa.PublicKey) error {
	bc.PoolMu.Lock()
	defer bc.PoolMu.Unlock()

	fmt.Printf("Adding transaction to pool: %+v\n", transaction)

	err := bc.ValidateTransaction(transaction, r, s, publicKey)
	if err != nil {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	bc.TransactionPool = append(bc.TransactionPool, transaction)
	bc.TxChan <- struct{}{}

	return nil
}

func (bc *Blockchain) ValidateTransaction(transaction PatientRecord, r, s *big.Int, publicKey ecdsa.PublicKey) error {
	hash := sha256.Sum256([]byte(transaction.PersonalData.ID))
	valid := ecdsa.Verify(&publicKey, hash[:], r, s)
	if !valid {
		return errors.New("invalid transaction signature")
	}
	return nil
}

func (bc *Blockchain) ProcessTransactions(passphrase string) {
	for {
		fmt.Println("Waiting for a new transaction signal...")
		<-bc.TxChan
		fmt.Println("Received a new transaction signal")
		bc.PoolMu.Lock()

		if len(bc.TransactionPool) == 0 {
			bc.PoolMu.Unlock()
			time.Sleep(time.Second)
			continue
		}

		transaction := bc.TransactionPool[0]
		bc.TransactionPool = bc.TransactionPool[1:]
		bc.PoolMu.Unlock()

		for _, newRecord := range transaction.MedicalRecords {
			bc.AddMedicalRecord(transaction.PersonalData.ID, newRecord, passphrase)
		}

		select {
		case bc.TxChan <- struct{}{}:
		default:
		}
	}
}
