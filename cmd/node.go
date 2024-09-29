package cmd

import (
	"crypto/ed25519"
	"fmt"

	"github.com/spf13/cobra"
)

var authorityAddress string

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start a node",
	Long:  `Start a node either as an authority node or as a client node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if authorityAddress == "" {
			authorityPrivateKey, authorityPublicKey, _ := ed25519.GenerateKey(nil)
			nodeInstance := NewNode(authorityPrivateKey, authorityPublicKey, "localhost:8080")
			authorityNode := NewAuthorityNode(authorityPrivateKey, nodeInstance)
			fmt.Println("Starting Authority Node...")
		} else {
			clientPrivateKey, clientPublicKey, _ := ed25519.GenerateKey(nil)
			nodeInstance := NewNode(clientPrivateKey, clientPublicKey, authorityAddress)
			fmt.Printf("Starting Client Node... Connecting to Authority Node at %s\n", authorityAddress)
		}
	},
}

func init() {
	nodeCmd.Flags().StringVarP(&authorityAddress, "authority", "a", "", "IP address of the authority node")
	rootCmd.AddCommand(nodeCmd)
}
