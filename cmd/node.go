package cmd

import (
	"crypto/ed25519"
	"fmt"

	"github.com/spf13/cobra"
)

var authorityAddress string

// TODO: Port hinzufügen per Parameter -p --port
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start a node",
	Long:  `Start a node either as an authority node or as a client node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if authorityAddress == "" {
			authorityPublicKey, authorityPrivateKey, _ := ed25519.GenerateKey(nil)
			authorityNode := NewAuthorityNode(authorityPrivateKey, authorityPublicKey)
			fmt.Println("Starting Authority Node...")
      authorityNode.SetupAuthorityNodeRoutes()
      authorityNode.Listen(":8080")
		} else {
			clientPublicKey, clientPrivateKey, _ := ed25519.GenerateKey(nil)
			node := NewNode(clientPrivateKey, clientPublicKey, authorityAddress)
			fmt.Printf("Starting Client Node... Connecting to Authority Node at %s\n", authorityAddress)
      node.SetupNodeRoutes()
      go node.AuthorityNodeDiscovery()
      node.Listen(":8080")
		}
	},
}

func init() {
	nodeCmd.Flags().StringVarP(&authorityAddress, "authority", "a", "", "IP address of the authority node")
	rootCmd.AddCommand(nodeCmd)
}
