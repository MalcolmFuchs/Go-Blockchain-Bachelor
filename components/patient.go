package components

import "fmt"

func (bc *Blockchain) AddPatient(patient PersonalData) {
	bc.Mu.Lock()
	defer bc.Mu.Unlock()
	bc.Patients[patient.ID] = patient
}

func (bc *Blockchain) GetPatient(id string) *PersonalData {

	fmt.Println(bc.Patients)

	for _, patient := range bc.Patients {
		if patient.ID == id {
			return &patient
		}
	}
	return nil
}
