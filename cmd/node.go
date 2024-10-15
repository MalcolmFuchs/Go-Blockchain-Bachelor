package cmd

import (
	"fmt"
	"os"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
	"github.com/spf13/cobra"
)

var (
	authorityAddress string
	port             string
)

// TODO: Port hinzufügen per Parameter -p --port
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Start a node",
	Long:  `Start a node either as an authority node or as a client node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if authorityAddress == "" {
			authorityNodePrivateKey, authorityNodePublicKey, err := utils.LoadPrivateKey(privKeyFile)
			if err != nil {
				fmt.Println("Fehler beim Laden des privaten Schlüssels:", err)
				os.Exit(1)
			}
			authorityNode := NewAuthorityNode(authorityNodePrivateKey, authorityNodePublicKey)
			fmt.Println("Starting Authority Node...")
			authorityNode.SetupAuthorityNodeRoutes()
			authorityNode.Listen(":" + port)
		} else {
			node := NewNode(nil, authorityAddress)
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
