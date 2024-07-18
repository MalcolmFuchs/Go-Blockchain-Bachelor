package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

var blockchainInstance *blockchain.Blockchain
var passphrase = "mysecretphrase12mysecretphrase12"

func init() {
	blockchainInstance = blockchain.CreateBlockchain()

	nodeNames := []string{"AOK", "TK", "Barmenia"}
	nodes := []blockchain.AuthorityNode{}

	for i, name := range nodeNames {
		privateKey, publicKey := blockchain.GenerateKeyPair()
		node := blockchain.AuthorityNode{
			ID:         fmt.Sprintf("%d", i+1),
			Name:       name,
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		}
		nodes = append(nodes, node)
	}

	blockchainInstance.Nodes = nodes
}

func addPatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient blockchain.PersonalData
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := blockchain.HashInsuranceNumber(patient.InsuranceNumber)
	patient.ID = id

	blockchainInstance.AddPatient(patient)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

func getPatientHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing hashed insurance number", http.StatusBadRequest)
		return
	}

	patient := blockchainInstance.GetPatient(id)
	if patient == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(patient)
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
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	blockchainInstance.AddMedicalRecord(id, record, passphrase)
	json.NewEncoder(w).Encode(blockchainInstance.Blocks)
}

func getMedicalRecordsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	passphrase := r.URL.Query().Get("passphrase")
	if id == "" || passphrase == "" {
		http.Error(w, "Missing id or passphrase", http.StatusBadRequest)
		return
	}

	records := blockchainInstance.GetMedicalRecords(id, passphrase)
	if records == nil {
		http.Error(w, "No records found or access denied", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(records)
}

func main() {
	// Starte den HTTP-Server
	http.HandleFunc("/blockchain", getBlockchain)
	http.HandleFunc("/addRecord", addMedicalRecordHandler)
	http.HandleFunc("/getRecords", getMedicalRecordsHandler)
	http.HandleFunc("/addPatient", addPatientHandler)
	http.HandleFunc("/getPatient", getPatientHandler)
	http.ListenAndServe(":8080", nil)
}
