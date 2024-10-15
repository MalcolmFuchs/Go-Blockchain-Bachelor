package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
	"github.com/spf13/cobra"
)

var (
	viewNodeAddress string
	patientKeyFile  string
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Zeigt alle Transaktionen eines Patienten an",
	Run: func(cmd *cobra.Command, args []string) {
		// Lade den privaten Schlüssel des Patienten
		patientPrivKey, patientPubKey, err := utils.LoadPrivateKey(patientKeyFile)
		if err != nil {
			fmt.Println("Fehler beim Laden des privaten Schlüssels des Patienten:", err)
			os.Exit(1)
		}

		// Serialisiere den öffentlichen Schlüssel des Patienten und kodiere ihn in Base64
		serializedPubKey := utils.SerializePublicKey(patientPubKey)
		patientID := base64.URLEncoding.EncodeToString(serializedPubKey)

		// Baue die URL mit dem patientID-Parameter
		url := fmt.Sprintf("http://%s/getPatientTransactions?patientID=%s", viewNodeAddress, patientID)

		// Sende eine GET-Anfrage an den Authority Node
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Fehler beim Abrufen der Transaktionen:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("Fehlerhafte Antwort vom Server: %s\n", string(body))
			os.Exit(1)
		}

		// Lese die Antwort (die verschlüsselten Transaktionen)
		var transactions []*blockchain.Transaction
		err = json.NewDecoder(resp.Body).Decode(&transactions)
		if err != nil {
			fmt.Println("Fehler beim Dekodieren der Transaktionen:", err)
			os.Exit(1)
		}

		// Verarbeite und entschlüssle die Transaktionen
		for _, tx := range transactions {
			// Deserialisiere den öffentlichen Schlüssel des Arztes
			doctorPubKey, err := utils.DeserializePublicKey(tx.Doctor)
			if err != nil {
				fmt.Println("Fehler beim Deserialisieren des öffentlichen Schlüssels des Arztes:", err)
				continue
			}

			// Konvertiere die Schlüssel zu ECDH-Schlüsseln
			patientEcdhPrivKey, err := utils.EcdsaPrivToEcdh(patientPrivKey)
			if err != nil {
				fmt.Println("Fehler bei der Konvertierung des privaten Schlüssels des Patienten:", err)
				continue
			}

			doctorEcdhPubKey, err := utils.EcdsaPubToEcdh(doctorPubKey)
			if err != nil {
				fmt.Println("Fehler bei der Konvertierung des öffentlichen Schlüssels des Arztes:", err)
				continue
			}

			// Entschlüssele die Daten
			plaintext, err := utils.DecryptData(patientEcdhPrivKey, doctorEcdhPubKey, tx.EncryptedData.Ciphertext, tx.EncryptedData.Nonce)
			if err != nil {
				fmt.Println("Fehler beim Entschlüsseln der Transaktion:", err)
				continue
			}

			// Dekodiere die entschlüsselten Daten in TransactionData
			var txData blockchain.TransactionData
			err = json.Unmarshal(plaintext, &txData)
			if err != nil {
				fmt.Println("Fehler beim Dekodieren der Transaktionsdaten:", err)
				continue
			}

			// Zeige die entschlüsselten Transaktionsdaten an
			fmt.Println("-----------")
			fmt.Printf("Transaktions-Hash: %x\n", tx.Hash)
			fmt.Printf("Typ: %s\n", txData.Type)
			fmt.Printf("Notizen: %s\n", txData.Notes)
			fmt.Printf("Ergebnisse: %s\n", txData.Results)
			fmt.Println("-----------")
		}
	},
}

func init() {
	viewCmd.Flags().StringVarP(&viewNodeAddress, "node_address", "a", "localhost:8080", "Adresse des Authority Nodes")
	viewCmd.Flags().StringVarP(&patientKeyFile, "key", "k", "", "Pfad zum privaten Schlüssel des Patienten (erforderlich)")
	viewCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(viewCmd)
}
