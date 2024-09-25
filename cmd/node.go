package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	blockchain "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/block"
	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
	"github.com/spf13/cobra"
)

var authorityIP string

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		if authorityIP == "" {
			startAuthorityNode()
		} else {
			startClientNode(authorityIP)
		}
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.Flags().StringVarP(&authorityIP, "authority", "a", "", "Specify authority node IP address")
}

func startAuthorityNode() {
	fmt.Println("Starting Authority Node...")

	privateKey, _, err := utils.GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		os.Exit(1)
	}

	bc := &blockchain.Blockchain{
		Patients:        make(map[string]blockchain.PersonalData),
		TransactionPool: []blockchain.EncryptedPatientRecord{},
		Blocks:          []blockchain.Block{},
		TxChan:          make(chan struct{}),
		PrivateKey:      privateKey, // PrivateKey des Authority Nodes
	}

	http.HandleFunc("/addTransaction", func(w http.ResponseWriter, r *http.Request) {
		var transaction blockchain.EncryptedPatientRecord
		err := json.NewDecoder(r.Body).Decode(&transaction)
		if err != nil {
			http.Error(w, "Invalid transaction payload", http.StatusBadRequest)
			return
		}

		bc.AddEncryptedRecord(transaction)
		fmt.Println("Received encrypted transaction for patient:", transaction.PatientID)

		if len(bc.TransactionPool) >= 10 {
			bc.CreateBlock()
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Transaction added for patient %s", transaction.PatientID)
	})

	fmt.Println("Authority Node is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func startClientNode(authorityIP string) {
	encryptionKey := []byte("your-32-byte-encryption-key12345")

	fmt.Println("Starting Client Node...")

	record := blockchain.MedicalRecord{
		Date:     "2024-09-25",
		Type:     "Blood Test",
		Provider: "Health Clinic",
		Notes:    "All values normal",
		Results:  "Cholesterol: 180 mg/dL",
	}

	err := blockchain.SendMedicalRecord("12345", record, encryptionKey, fmt.Sprintf("http://%s:8080", authorityIP))
	if err != nil {
		fmt.Printf("Error sending medical record: %v\n", err)
	} else {
		fmt.Println("Medical record successfully sent.")
	}
}

func sendTransactionToAuthorityNode(transaction blockchain.EncryptedPatientRecord, authorityIP string) error {
	url := fmt.Sprintf("http://%s:8080/addTransaction", authorityIP)

	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("error serializing transaction: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending transaction to authority node: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from authority node: %v", resp.Status)
	}

	return nil
}
