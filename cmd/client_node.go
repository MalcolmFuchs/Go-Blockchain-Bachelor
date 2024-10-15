package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

type Node struct {
	Blockchain           *blockchain.Blockchain
	Doctors              map[string]DoctorData
	Patients             map[string]PatientData
	TrustedPublicKey     ed25519.PublicKey
	AuthorityNodeAddress string
	PrivateKey           ed25519.PrivateKey
}

func NewNode(privateKey ed25519.PrivateKey, trustedPublicKey ed25519.PublicKey, authorityNodeAddress string) *Node {
	return &Node{
		Blockchain:           &blockchain.Blockchain{Blocks: []*blockchain.Block{}, BlockMap: make(map[string]*blockchain.Block)},
		Doctors:              make(map[string]DoctorData),
		Patients:             make(map[string]PatientData),
		TrustedPublicKey:     trustedPublicKey,
		PrivateKey:           privateKey,
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
	PublicKey    ed25519.PublicKey                  `json:"public_key"`
}

func (n *Node) AddDoctor(doctor DoctorData) {
	doctorKeyHash := hex.EncodeToString(doctor.PublicKey)
	n.Doctors[doctorKeyHash] = doctor
}

func (n *Node) GetDoctor(publicKey ed25519.PublicKey) (*DoctorData, error) {
	doctorKeyHash := hex.EncodeToString(publicKey)

	if doctor, exists := n.Doctors[doctorKeyHash]; exists {
		return &doctor, nil
	}
	return nil, fmt.Errorf("doctor not found")
}

func (n *Node) AddPatient(patient PatientData) {
	patientKeyHash := hex.EncodeToString(patient.PublicKey)
	n.Patients[patientKeyHash] = patient
}

func (n *Node) GetPatient(publicKey ed25519.PublicKey) (*PatientData, error) {
	patientKeyHash := hex.EncodeToString(publicKey)
	if patient, exists := n.Patients[patientKeyHash]; exists {
		return &patient, nil
	}
	return nil, fmt.Errorf("patient not found")
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

	publicKeyHex, ok := result["publicKey"]
	if !ok {
		fmt.Println("Public Key nicht in der Antwort gefunden")
		return
	}

	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		fmt.Printf("Fehler beim Dekodieren des Public Keys: %v\n", err)
		return
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		fmt.Println("Ungültige Länge des Public Keys")
		return
	}

	n.TrustedPublicKey = ed25519.PublicKey(publicKeyBytes)

	fmt.Println("Authority Node Public Key erfolgreich erhalten und gespeichert")

	// Sync Blockchain
	err = n.SyncWithAuthorityNode(n.AuthorityNodeAddress)
	if err != nil {
		fmt.Printf("Fehler bei der Synchronisierung mit dem Authority Node: %v\n", err)
		return
	}

	fmt.Println("Blockchain erfolgreich synchronisiert")
}
