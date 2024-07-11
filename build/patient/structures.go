package patient

import "time"

type PatientRecord struct {
	ID             string          `json:"id"`
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
