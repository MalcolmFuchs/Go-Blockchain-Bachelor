package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	blockchain "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/block"
	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

type TransactionRequest struct {
	PatientData  string `json:"patientData"`
	DoctorPubKey string `json:"doctorPubKey"` // Public Key des Arztes in Base64-Form
}

func StartClientNode(authorityIP string) {
	http.HandleFunc("/addTransaction", handleTransactionRequest)
	fmt.Println("Client Node running. Send transactions via POST to /addTransaction")
	http.ListenAndServe(":8080", nil) // Client Node API läuft auf Port 8080
}

func handleTransactionRequest(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Public Key des Arztes dekodieren
	doctorPublicKey, err := utils.DecodePublicKey(req.DoctorPubKey)
	if err != nil {
		http.Error(w, "Invalid doctor public key", http.StatusBadRequest)
		return
	}

	aesKey := utils.GenerateAESKey()

	// Patientendaten mit AES-GCM verschlüsseln
	encryptedData, err := utils.EncryptData(aesKey, []byte(req.PatientData))
	if err != nil {
		http.Error(w, "Error encrypting patient data", http.StatusInternalServerError)
		return
	}

	patientPrivateKey := utils.GetPatientPrivateKey()
	patientSign := utils.SignTransaction(patientPrivateKey, encryptedData)

	// AES-Schlüssel mit dem Public Key des Arztes verschlüsseln
	encryptedAESKey, err := utils.EncryptAESKeyWithDoctorKey(doctorPublicKey, aesKey, patientPrivateKey)
	if err != nil {
		http.Error(w, "Error encrypting AES key", http.StatusInternalServerError)
		return
	}

	transaction := blockchain.NewTransaction(encryptedData, encryptedAESKey, patientSign)

	SendTransactionToAuthority(transaction)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Transaction successfully created and sent")
}

func SendTransactionToAuthority(transaction blockchain.Transaction) {
	conn, err := net.Dial("tcp", "localhost:8081") // Sende an Authority Node
	if err != nil {
		fmt.Println("Error connecting to authority node:", err)
		return
	}
	defer conn.Close()

	err = json.NewEncoder(conn).Encode(transaction)
	if err != nil {
		fmt.Println("Error sending transaction:", err)
		return
	}
	fmt.Println("Transaction sent to authority node successfully")
}
