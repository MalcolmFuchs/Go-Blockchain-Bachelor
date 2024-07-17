package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

var blockchainInstance *blockchain.Blockchain
var passphrase = "mysecretphrase12mysecretphrase12"

func init() {
	blockchainInstance = blockchain.CreateBlockchain()

	privateKey1, publicKey1 := GenerateKeyPair()
	privateKey2, publicKey2 := GenerateKeyPair()
	privateKey3, publicKey3 := GenerateKeyPair()
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
			PrivateKey: privateKey3,
			PublicKey:  publicKey3,
		},
	}

	blockchainInstance.Nodes = nodes
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchainInstance.Blocks)
}

// TO-DO:
// Query muss verschlüsselt werden, InsuranceNumber darf nicht sichtbar sein!!

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

func addPatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient blockchain.PersonalData
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blockchainInstance.AddPatient(patient)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

func getPatientHandler(w http.ResponseWriter, r *http.Request) {
	insuranceNumber := r.URL.Query().Get("insuranceNumber")
	if insuranceNumber == "" {
		http.Error(w, "Missing insurance number", http.StatusBadRequest)
		return
	}

	patient := blockchainInstance.GetPatient(insuranceNumber)
	if patient == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(patient)
}

func main() {
	// Initial Setup
	fmt.Println("Blockchain initialized with nodes:")
	for _, node := range blockchainInstance.Nodes {
		fmt.Printf("Node Name: %s\n", node.Name)
	}

	// Starte den HTTP-Server
	http.HandleFunc("/blockchain", getBlockchain)
	http.HandleFunc("/addRecord", addMedicalRecordHandler)
	http.HandleFunc("/getRecords", getMedicalRecordsHandler)
	http.HandleFunc("/addPatient", addPatientHandler)
	http.HandleFunc("/getPatient", getPatientHandler)
	http.ListenAndServe(":8080", nil)
}

func GenerateKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
	}
	return privateKey, privateKey.PublicKey
}
