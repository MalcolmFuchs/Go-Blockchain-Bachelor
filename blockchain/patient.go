package blockchain

func (bc *Blockchain) AddPatient(patient PersonalData) {
	bc.Patients[patient.InsuranceNumber] = patient
}

func (bc *Blockchain) GetPatient(insuranceNumber string) *PersonalData {
	patient, found := bc.Patients[insuranceNumber]
	if !found {
		return nil
	}
	return &patient
}
