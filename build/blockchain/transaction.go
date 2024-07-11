package blockchain

import (
	patient "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/build/patient"
)

func NewTransaction(patientID string, record patient.PatientRecord) Transaction {
	return Transaction{PatientID: patientID, Record: record}
}
