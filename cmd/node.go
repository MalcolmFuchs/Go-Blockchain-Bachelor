package cmd

import (
	"crypto/ed25519"
	"fmt"

	"github.com/spf13/cobra"
)

// Definition der Variablen für Flags
var authorityAddress string

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start a node",
	Long:  `Start a node either as an authority node or as a client `,
	Run: func(cmd *cobra.Command, args []string) {
		if authorityAddress == "" {
			// Authority Node erstellen
			authorityPrivateKey, authorityPublicKey, _ := ed25519.GenerateKey(nil)
			nodeInstance := NewNode(authorityPrivateKey, authorityPublicKey, "localhost:8080")
			authorityNode := NewAuthorityNode(authorityPrivateKey, nodeInstance)
			fmt.Println("Starting Authority ..")
			// Starte die Blockerstellungslogik oder API
		} else {
			// Client Node erstellen
			clientPrivateKey, clientPublicKey, _ := ed25519.GenerateKey(nil)
			nodeInstance := NewNode(clientPrivateKey, clientPublicKey, authorityAddress)
			fmt.Printf("Starting Client Connecting to Authority Node at %s\n", authorityAddress)
			// Logik für den Client Node (z.B. Synchronisation)
		}
	},
}

func init() {
	// Flag für die IP-Adresse des Authority Nodes hinzufügen
	nodeCmd.Flags().StringVarP(&authorityAddress, "authority", "a", "", "IP address of the authority node")
	rootCmd.AddCommand(nodeCmd)
}
