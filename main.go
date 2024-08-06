package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/components"
	blockchain "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/components"
)

var blockchainInstance *blockchain.Blockchain
var passphrase = "mysecretphrase12mysecretphrase12"

func init() {
	blockchainInstance = blockchain.CreateBlockchain()
	go blockchainInstance.ProcessTransactions(passphrase)

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
	var patient components.PersonalData
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	patient.ID = fmt.Sprintf("%x", sha256.Sum256([]byte(patient.InsuranceNumber)))

	hash := sha256.Sum256([]byte(patient.ID))
	rSign, sSign, err := ecdsa.Sign(rand.Reader, blockchainInstance.Nodes[0].PrivateKey, hash[:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transaction := components.PatientRecord{
		PersonalData: patient,
	}

	err = blockchainInstance.AddTransactionToPool(transaction, rSign, sSign, blockchainInstance.Nodes[0].PublicKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blockchainInstance.Mu.Lock()
	blockchainInstance.Patients[patient.ID] = patient
	blockchainInstance.Mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

func getPatientHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	blockchainInstance.Mu.Lock()
	defer blockchainInstance.Mu.Unlock()

	patient, exists := blockchainInstance.Patients[id]
	if !exists {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(patient)
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	blockchainInstance.Mu.Lock()
	defer blockchainInstance.Mu.Unlock()

	json.NewEncoder(w).Encode(blockchainInstance.Blocks)
}

func addMedicalRecordHandler(w http.ResponseWriter, r *http.Request) {
	var record components.MedicalRecord
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

	blockchainInstance.Mu.Lock()
	defer blockchainInstance.Mu.Unlock()

	patient, exists := blockchainInstance.Patients[id]
	if !exists {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	var patientRecord components.PatientRecord
	for i, record := range blockchainInstance.TransactionPool {
		if record.PersonalData.ID == id {
			patientRecord = blockchainInstance.TransactionPool[i]
			blockchainInstance.TransactionPool = append(blockchainInstance.TransactionPool[:i], blockchainInstance.TransactionPool[i+1:]...)
			break
		}
	}

	if patientRecord.PersonalData.ID == "" {
		patientRecord = components.PatientRecord{
			PersonalData: patient,
		}
	}

	record.Date = time.Now()
	patientRecord.MedicalRecords = append(patientRecord.MedicalRecords, record)

	hash := sha256.Sum256([]byte(patient.ID))
	rSign, sSign, err := ecdsa.Sign(rand.Reader, blockchainInstance.Nodes[0].PrivateKey, hash[:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = blockchainInstance.AddTransactionToPool(patientRecord, rSign, sSign, blockchainInstance.Nodes[0].PublicKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(blockchainInstance.Blocks)
}

func getMedicalRecordsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	passphrase := r.URL.Query().Get("passphrase")
	if id == "" || passphrase == "" {
		http.Error(w, "Missing ID or passphrase", http.StatusBadRequest)
		return
	}

	records, err := blockchainInstance.GetMedicalRecords(id, passphrase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// if records == nil {
	// 	http.Error(w, "No records found or access denied", http.StatusForbidden)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(records)
}

func main() {
	http.HandleFunc("/blockchain", getBlockchain)
	http.HandleFunc("/addRecord", addMedicalRecordHandler)
	http.HandleFunc("/getRecords", getMedicalRecordsHandler)
	http.HandleFunc("/addPatient", addPatientHandler)
	http.HandleFunc("/getPatient", getPatientHandler)

	fmt.Println("Server listening on port :8080")
	http.ListenAndServe(":8080", nil)
}
