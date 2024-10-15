package cmd

import (
	"crypto/ed25519"
	"fmt"

	"github.com/spf13/cobra"
)

var authorityAddress string
var port string

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
			authorityNode.Listen(":" + port)
		} else {
			_, clientPrivateKey, _ := ed25519.GenerateKey(nil)
			node := NewNode(clientPrivateKey, nil, authorityAddress)
			fmt.Printf("Starting Client Node... Connecting to Authority Node at %s\n", authorityAddress)
			node.SetupNodeRoutes()
      go node.StartSyncRoutine()
			node.Listen(":" + port)
		}
	},
}

func init() {
	nodeCmd.Flags().StringVarP(&authorityAddress, "authority", "a", "", "IP address of the authority node")
	nodeCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port für den Node")
	rootCmd.AddCommand(nodeCmd)
}
