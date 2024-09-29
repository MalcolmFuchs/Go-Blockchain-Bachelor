package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net/http"

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

func (a *Node) Listen(addr string) {
  http.ListenAndServe(addr, nil)
}

// TODO: Build discovery function ()
func (a *Node) AuthorityNodeDiscovery() {}
