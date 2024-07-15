package main

import (
	"fmt"
	"time"

	bc "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/build/blockchain"
	patient "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/build/patient"
)

func main() {
	nodes := []bc.AuthorityNode{
		{Id: "Gesundheitsministerium", PublicKey: "key1", PrivateKey: "key1"},
		{Id: "AOK", PublicKey: "key2", PrivateKey: "key2"},
		{Id: "TK", PublicKey: "key3", PrivateKey: "key3"},
	}

	blockchain := bc.CreateBlockchain(nodes)

	// Beispiel für die Erstellung eines PatientRecord.
	patient1 := patient.PatientRecord{
		ID: "1234567890",
		PersonalData: patient.PersonalData{
			FirstName:       "Max",
			LastName:        "Mustermann",
			BirthDate:       time.Date(1998, 9, 23, 0, 0, 0, 0, time.UTC),
			InsuranceNumber: "AOK123456789",
		},
		MedicalRecords: []patient.MedicalRecord{
			{
				Date:     time.Date(2024, 7, 11, 12, 46, 55, 0, time.UTC),
				Type:     "Arztbrief",
				Provider: "Dr. Med. Mann",
				Notes:    "Zeigt Symptomatik von Fieber, Husten und Halsschmerzen",
				Results:  "Patient leidet an Grippe",
			},
			{
				Date:     time.Date(2024, 7, 11, 13, 26, 15, 0, time.UTC),
				Type:     "Medikationspläne",
				Provider: "Rosenapotheke",
				Notes:    "Antibiotikum, 3x täglich einnehmen",
				Results:  "Amoxihexal",
			},
		},
	}

	transaction := bc.NewTransaction("1", patient1)

	// Erstellen eines neuen Blocks
	newBlock := nodes[0].CreateBlock(
		[]bc.Transaction{transaction},
		&blockchain.Chain[len(blockchain.Chain)-1],
		nodes,
		&blockchain,
	)

	// Überprüfen ob der neue Block gültig ist
	if newBlock != nil && nodes[0].ValidateBlock(newBlock, &blockchain) {
		// Hinzufügen des neuen Blocks zur Blockchain
		blockchain.AddBlock(*newBlock)
	} else {
		fmt.Println("Der neue Block ist ungültig")
	}

	// Überprüfung ob die Blockchain gültig ist
	if !blockchain.IsValid() {
		fmt.Println("Die Blockchain ist ungültig")
	}

	// Blockchain anzeigen
	for i, block := range blockchain.Chain {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("\tTimestamp: %s\n", block.Timestamp)
		fmt.Printf("\tPrev. Hash: %x\n", block.PrevHash)
		fmt.Printf("\tHash: %x\n", block.Hash)
		fmt.Println("\tTransactions:")
		for _, tx := range block.Transactions {
			fmt.Printf("\t\tPatient ID: %s\n", tx.PatientID)
			fmt.Printf("\t\tRecord: %+v\n", tx.Record)
		}
	}

	// Hash des PatientRecord erzeugen.
	// hash := patient1.PatientHash()
	// println("Hash des PatientRecord: ", hash)
}
