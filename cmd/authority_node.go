package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"crypto/ed25519"

	blockchain "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/block"
	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

// StartAuthorityNode startet den Authority Node und verwaltet Transaktionen sowie Blockerstellung
func StartAuthorityNode(authorityPrivateKey ed25519.PrivateKey) {
	ticker := time.NewTicker(5 * time.Minute)
	transactionPool := make([]blockchain.Transaction, 0)

	for {
		select {
		case <-ticker.C:
			if len(transactionPool) > 0 {
				createAndSignBlock(transactionPool, authorityPrivateKey)
				transactionPool = []blockchain.Transaction{}
			}
		default:
			newTransaction, err := ReceiveTransactionFromClient()
			if err != nil {
				fmt.Println("Error receiving transaction:", err)
				continue
			}

			patientPublicKey := utils.GetPatientPublicKey(newTransaction.PatientSign)
			if !utils.VerifySignature(patientPublicKey, newTransaction.EncryptedPatientData, newTransaction.PatientSign) {
				fmt.Println("Invalid patient signature")
				continue
			}

			if newTransaction.DoctorSign != nil {
				doctorPublicKey := utils.GetDoctorPublicKey()
				if !utils.VerifySignature(doctorPublicKey, newTransaction.EncryptedPatientData, newTransaction.DoctorSign) {
					fmt.Println("Invalid doctor signature")
					continue
				}
			}

			transactionPool = append(transactionPool, newTransaction)

			if len(transactionPool) >= 10 {
				createAndSignBlock(transactionPool, authorityPrivateKey)
				transactionPool = []blockchain.Transaction{}
			}
		}
	}
}

func createAndSignBlock(transactions []blockchain.Transaction, authorityPrivateKey ed25519.PrivateKey) {
	lastBlock := blockchain.GetLastBlock()

	newBlock := blockchain.NewBlock(transactions, lastBlock.Hash)

	utils.SignBlock(&newBlock, authorityPrivateKey)

	blockchain.AddBlock(newBlock)

	fmt.Println("Block created and signed with", len(transactions), "transactions")
}

// ReceiveTransactionFromClient empfängt eine Transaktion von einem Client Node über TCP
func ReceiveTransactionFromClient() (blockchain.Transaction, error) {
	ln, err := net.Listen("tcp", ":8081") // Authority Node lauscht auf Port 8081
	if err != nil {
		return blockchain.Transaction{}, fmt.Errorf("Error setting up server: %v", err)
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		return blockchain.Transaction{}, fmt.Errorf("Error accepting connection: %v", err)
	}
	defer conn.Close()

	var transaction blockchain.Transaction
	err = json.NewDecoder(conn).Decode(&transaction)
	if err != nil {
		return blockchain.Transaction{}, fmt.Errorf("Error decoding transaction: %v", err)
	}

	return transaction, nil
}
