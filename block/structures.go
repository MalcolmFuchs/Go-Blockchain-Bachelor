package blockchain

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
	Index       int                      `json:"index"`
	Timestamp   time.Time                `json:"timestamp"`
	PatientData []EncryptedPatientRecord `json:"patientData"`
	Hash        string                   `json:"hash"`
	PrevHash    string                   `json:"prevHash"`
	SignatureR  string                   `json:"signatureR"`
	SignatureS  string                   `json:"signatureS"`
}

type Blockchain struct {
	Blocks          []Block                  `json:"blocks"`
	Patients        map[string]PersonalData  `json:"patients"`
	TransactionPool []EncryptedPatientRecord `json:"transactionPool"`
	Mu              sync.Mutex
	PoolMu          sync.Mutex
	TxChan          chan struct{}
	PrivateKey      *ecdsa.PrivateKey
}

type PatientRecord struct {
	PersonalData   PersonalData    `json:"personalData"`
	MedicalRecords []MedicalRecord `json:"medicalRecords"`
}

type EncryptedPatientRecord struct {
	PatientID       string                 `json:"patientID"`
	EncryptedRecord EncryptedMedicalRecord `json:"encryptedRecord"`
}

type PersonalData struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	BirthDate       time.Time `json:"birthDate"`
	InsuranceNumber string    `json:"insuranceNumber"`
}

type EncryptedPersonalData struct {
	FirstName       string                   `json:"firstName"`
	LastName        string                   `json:"lastName"`
	BirthDate       string                   `json:"birthDate"`
	InsuranceNumber string                   `json:"insuranceNumber"`
	MedicalRecords  []EncryptedMedicalRecord `json:"medicalRecords"`
}

type MedicalRecord struct {
	Date     string `json:"date"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
	Notes    string `json:"notes"`
	Results  string `json:"results"`
}

type MedicalRecordTransaction struct {
	PatientID       string                 `json:"patientID"`
	EncryptedRecord EncryptedMedicalRecord `json:"encryptedRecord"`
}

type EncryptedMedicalRecord struct {
	Date     string `json:"date"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
	Notes    string `json:"notes"`
	Results  string `json:"results"`
}
