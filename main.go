package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

var blockchainInstance *blockchain.Blockchain
var passphrase = "mysecretphrase123"

func init() {
	blockchainInstance = blockchain.CreateBlockchain()

	privateKey1, publicKey1 := GenerateKeyPair()
	privateKey2, publicKey2 := GenerateKeyPair()
	nodes := []blockchain.AuthorityNode{
		{
			ID:         "1",
			Name:       "AOK",
			PrivateKey: privateKey1,
			PublicKey:  publicKey1,
		},
		{
			ID:         "2",
			Name:       "TK",
			PrivateKey: privateKey2,
			PublicKey:  publicKey2,
		},
		{
			ID:         "3",
			Name:       "Barmenia",
			PrivateKey: privateKey2,
			PublicKey:  publicKey2,
		},
	}

	blockchainInstance.Nodes = nodes
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchainInstance.Blocks)
}

func addMedicalRecordHandler(w http.ResponseWriter, r *http.Request) {
	var record blockchain.MedicalRecord
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	insuranceNumber := r.URL.Query().Get("insuranceNumber")
	if insuranceNumber == "" {
		http.Error(w, "Missing insurance number", http.StatusBadRequest)
		return
	}

	blockchainInstance.AddMedicalRecord(insuranceNumber, record, passphrase)
	json.NewEncoder(w).Encode(blockchainInstance.Blocks)
}

func getMedicalRecordsHandler(w http.ResponseWriter, r *http.Request) {
	insuranceNumber := r.URL.Query().Get("insuranceNumber")
	passphrase := r.URL.Query().Get("passphrase")
	if insuranceNumber == "" || passphrase == "" {
		http.Error(w, "Missing insurance number or passphrase", http.StatusBadRequest)
		return
	}

	records := blockchainInstance.GetMedicalRecords(insuranceNumber, passphrase)
	if records == nil {
		http.Error(w, "No records found or access denied", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(records)
}

func main() {
	// Initial Setup
	fmt.Println("Blockchain initialized with nodes:")
	for _, node := range blockchainInstance.Nodes {
		fmt.Printf("Node Name: %s\n", node.Name)
	}

	// Patient anlegen und erste Daten hinzufügen
	newRecord1 := blockchain.MedicalRecord{
		Date:     time.Now(),
		Type:     "Checkup",
		Provider: "Dr. Smith",
		Notes:    "Patient in good health.",
		Results:  "All tests normal.",
	}
	blockchainInstance.AddMedicalRecord("1234567890", newRecord1, passphrase)

	// Weitere Daten hinzufügen
	newRecord2 := blockchain.MedicalRecord{
		Date:     time.Now().AddDate(0, 1, 0), // 1 Monat später
		Type:     "Blood Test",
		Provider: "LabCorp",
		Notes:    "Cholesterol level slightly high.",
		Results:  "Cholesterol: 210 mg/dL",
	}
	blockchainInstance.AddMedicalRecord("1234567890", newRecord2, passphrase)

	// Daten anzeigen mit erlaubtem Zugriff
	fmt.Println("Retrieving records with correct passphrase:")
	records := blockchainInstance.GetMedicalRecords("1234567890", passphrase)
	if records != nil {
		for _, record := range records {
			fmt.Printf("Date: %s\n", record.Date)
			fmt.Printf("Type: %s\n", record.Type)
			fmt.Printf("Provider: %s\n", record.Provider)
			fmt.Printf("Notes: %s\n", record.Notes)
			fmt.Printf("Results: %s\n", record.Results)
			fmt.Println()
		}
	} else {
		fmt.Println("No records found or access denied")
	}

	// Daten anzeigen mit verweigertem Zugriff
	fmt.Println("Retrieving records with incorrect passphrase:")
	records = blockchainInstance.GetMedicalRecords("1234567890", "wrongpassphrase")
	if records != nil {
		for _, record := range records {
			fmt.Printf("Date: %s\n", record.Date)
			fmt.Printf("Type: %s\n", record.Type)
			fmt.Printf("Provider: %s\n", record.Provider)
			fmt.Printf("Notes: %s\n", record.Notes)
			fmt.Printf("Results: %s\n", record.Results)
			fmt.Println()
		}
	} else {
		fmt.Println("No records found or access denied")
	}

	// Starte den HTTP-Server
}

func GenerateKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
	}
	return privateKey, privateKey.PublicKey
}
