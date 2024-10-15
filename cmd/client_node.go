package cmd

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

type Node struct {
	Blockchain           *blockchain.Blockchain
	Doctors              map[string]DoctorData
	Patients             map[string]PatientData
	AuthorityNodeAddress string
	TrustedPublicKey     *ecdsa.PublicKey
}

func NewNode(privateKey *ecdsa.PrivateKey, authorityNodeAddress string) *Node {
	return &Node{
		Blockchain:           &blockchain.Blockchain{Blocks: []*blockchain.Block{}, BlockMap: make(map[string]*blockchain.Block)},
		Doctors:              make(map[string]DoctorData),
		Patients:             make(map[string]PatientData),
		AuthorityNodeAddress: authorityNodeAddress,
	}
}

type DoctorData struct {
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	PublicKey ed25519.PublicKey `json:"public_key"`
}

type PatientData struct {
	Transactions map[string]*blockchain.Transaction `json:"transactions"`
}

// Transaction to Authority-Client
func (n *Node) ForwardTransaction(transaction *blockchain.Transaction) error {
	if transaction == nil {
		return fmt.Errorf("transaction is nil")
	}

	fmt.Printf("Forwarding transaction with hash %x to authority node at %s\n", transaction.Hash, n.AuthorityNodeAddress)
	return nil
}

func (n *Node) Listen(addr string) {
	http.ListenAndServe(addr, nil)
}

// check if conditions are met every 5th minute with sleep
func (n *Node) StartSyncRoutine() {
	// Create the ticker inside the Go routine
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Sync Blockchain
		n.AuthorityNodeDiscovery()
	}
}

// Give the ClientNOde the publicKey of AuthorityNode
func (n *Node) AuthorityNodeDiscovery() {
	resp, err := http.Get(fmt.Sprintf("http://%s/getPublicKey", n.AuthorityNodeAddress))
	if err != nil {
		fmt.Printf("Fehler beim Verbinden mit dem Authority Node: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Fehlerhafte Antwort vom Authority Node: %d\n", resp.StatusCode)
		return
	}

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Fehler beim Dekodieren der Antwort: %v\n", err)
		return
	}

	publicKeyBytes, ok := result["publicKey"]
	if !ok {
		fmt.Println("Public Key nicht in der Antwort gefunden")
		return
	}

	publicKeyBytesDecoded, err := base64.StdEncoding.DecodeString(publicKeyBytes)
	if err != nil {
		fmt.Printf("error decoding base64 public key: %v", err)
		return
	}

	publicKey, err := utils.DeserializePublicKey(publicKeyBytesDecoded)
	if err != nil {
		fmt.Printf("couldn't deserialize doctor public key: %v", err)
		return
	}

	n.TrustedPublicKey = publicKey

	fmt.Println("Authority Node Public Key erfolgreich erhalten und gespeichert")

	// Sync Blockchain
	err = n.SyncWithAuthorityNode(n.AuthorityNodeAddress)
	if err != nil {
		fmt.Printf("Fehler bei der Synchronisierung mit dem Authority Node: %v\n", err)
		return
	}

	fmt.Println("Blockchain erfolgreich synchronisiert")
}
