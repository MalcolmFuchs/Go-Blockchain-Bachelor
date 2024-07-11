package patient

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Erzeugt einen SHA-256 Hash für PatientRecord.
func (p *PatientRecord) PatientHash() string {
	record, _ := json.Marshal(p)
	hash := sha256.Sum256(record)
	return hex.EncodeToString(hash[:])
}
