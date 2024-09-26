package cmd

import (
	"crypto/ed25519"
	"fmt"

	"github.com/spf13/cobra"
)

var authorityIP string
var authorityPrivateKey ed25519.PrivateKey

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.Flags().StringVarP(&authorityIP, "authority", "a", "", "IP address of the authority node")
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		if authorityIP != "" {
			StartClientNode(authorityIP)
		} else {
			fmt.Println("Starting Authority Node...")
			authorityPrivateKey = getAuthorityPrivateKey()
			StartAuthorityNode(authorityPrivateKey)
		}
	},
}

func getAuthorityPrivateKey() ed25519.PrivateKey {
	_, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic("Error generating authority private key")
	}
	return privateKey
}
