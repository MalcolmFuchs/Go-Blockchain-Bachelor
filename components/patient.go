package components

func (bc *Blockchain) AddPatient(patient PersonalData) {
	bc.Patients[patient.ID] = patient
}

func (bc *Blockchain) GetPatient(id string) *PersonalData {
	for _, patient := range bc.Patients {
		if patient.ID == id {
			return &patient
		}
	}
	return nil
}
