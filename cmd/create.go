package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
	"github.com/spf13/cobra"
)

var (
	nodeAddress string
	txType      string
	notes       string
	results     string
	pubKeyFile  string
	privKeyFile string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Erstellt eine neue Transaktion",
	Long:  "Dieser Befehl ermöglicht es, eine neue Transaktion lokal zu erstellen.",
	Run: func(cmd *cobra.Command, args []string) {
		// Lade den privaten Schlüssel des Arztes
		senderPrivKey, _, err := utils.LoadPrivateKey(privKeyFile)
		if err != nil {
			fmt.Println("Fehler beim Laden des privaten Schlüssels:", err)
			os.Exit(1)
		}

		patientPubKey, err := utils.LoadPublicKey(pubKeyFile)

		// Erstelle die Transaktion
		transaction, err := blockchain.NewTransaction(txType, notes, results, senderPrivKey, patientPubKey)
		if err != nil {
			fmt.Println("Fehler beim Erstellen der Transaktion:", err)
			os.Exit(1)
		}

		// Ausgabe der Transaktion
		txJSON, err := json.MarshalIndent(transaction, "", "  ")
		if err != nil {
			fmt.Println("Fehler beim Serialisieren der Transaktion:", err)
			os.Exit(1)
		}

		fmt.Println("Transaktion erfolgreich erstellt:")
		fmt.Println(string(txJSON))

		resp, err := http.Post(fmt.Sprintf("http://%s/addTransaction", nodeAddress), "application/json", bytes.NewBuffer(txJSON))
		if err != nil {
			fmt.Printf("failed to send sync request: %v", err)
      return
		}
		defer resp.Body.Close()
	},
}

func init() {
	createCmd.Flags().StringVarP(&nodeAddress, "node_address", "a", "", "Typ der Transaktion (erforderlich)")
	createCmd.Flags().StringVarP(&txType, "type", "t", "", "Typ der Transaktion (erforderlich)")
	createCmd.Flags().StringVarP(&notes, "notes", "n", "", "Notizen zur Transaktion")
	createCmd.Flags().StringVarP(&results, "results", "r", "", "Ergebnisse der Transaktion")
	createCmd.Flags().StringVarP(&pubKeyFile, "patient", "p", "", "Public Key des Patienten in Hex (erforderlich)")
	createCmd.Flags().StringVarP(&privKeyFile, "key", "k", "private_key.pem", "Pfad zum privaten Schlüssel des Arztes")

	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("doctor")
	createCmd.MarkFlagRequired("patient")

	rootCmd.AddCommand(createCmd)
}
