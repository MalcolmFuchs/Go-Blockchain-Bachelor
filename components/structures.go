package components

import (
	"crypto/ecdsa"
	"sync"
	"time"
)

type AuthorityNode struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	PrivateKey *ecdsa.PrivateKey `json:"-"`
	PublicKey  ecdsa.PublicKey   `json:"-"`
}

type Block struct {
	Index       int
	Timestamp   time.Time
	PatientData PatientRecord
	Hash        string
	PrevHash    string
	SignatureR  string
	SignatureS  string
}

type Blockchain struct {
	Blocks          []Block                 `json:"blocks"`
	Nodes           []AuthorityNode         `json:"nodes"`
	Patients        map[string]PersonalData `json:"patients"`
	TransactionPool []PatientRecord         `json:"transactionPool"`
	LastNodeIndex   int                     `json:"lastNodeIndex"`
	Mu              sync.Mutex
	PoolMu          sync.Mutex
	TxChan          chan struct{}
}

type PatientRecord struct {
	PersonalData   PersonalData    `json:"personalData"`
	MedicalRecords []MedicalRecord `json:"medicalRecords"`
}

// Persönliche Informationen eines Patienten.
type PersonalData struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	BirthDate       time.Time `json:"birthDate"`
	InsuranceNumber string    `json:"insuranceNumber"`
}

type EncryptedPersonalData struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	BirthDate       string `json:"birthDate"`
	InsuranceNumber string `json:"insuranceNumber"`
}

// Definiert die Struktur eines medizinischen Eintrags.
type MedicalRecord struct {
	Date     time.Time `json:"date"`
	Type     string    `json:"type"`
	Provider string    `json:"provider"`
	Notes    string    `json:"notes"`
	Results  string    `json:"results"`
}

type EncryptedMedicalRecord struct {
	Date     string `json:"date"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
	Notes    string `json:"notes"`
	Results  string `json:"results"`
}
