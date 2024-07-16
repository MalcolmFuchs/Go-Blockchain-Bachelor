package blockchain

import (
	"crypto/ecdsa"
	"time"
)

type AuthorityNode struct {
	ID         string
	Name       string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

type Block struct {
	Index       int
	Timestamp   string
	PatientData PatientRecord
	Hash        string
	PrevHash    string
	SignatureR  string
	SignatureS  string
}

type Blockchain struct {
	Blocks []Block
	Nodes  []AuthorityNode
}

type HealthRecord struct {
	PatientID string
	Data      string
}

type PatientRecord struct {
	PersonalData   PersonalData    `json:"personalData"`
	MedicalRecords []MedicalRecord `json:"medicalRecords"`
}

// Persönliche Informationen eines Patienten.
type PersonalData struct {
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	BirthDate       time.Time `json:"birthDate"`
	InsuranceNumber string    `json:"insuranceNumber"`
}

// Definiert die Struktur eines medizinischen Eintrags.
type MedicalRecord struct {
	Date     time.Time `json:"date"`
	Type     string    `json:"type"`
	Provider string    `json:"provider"`
	Notes    string    `json:"notes"`
	Results  string    `json:"results"`
}
